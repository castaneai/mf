package mf

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent         = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.73"
	sessionCookieName = "_moneybook_session"
)

type Client struct {
	hc   *http.Client
	opts *ClientOption
}

type ClientOption struct {
	Host      string
	SessionID string
}

func NewClient(hc *http.Client, opts *ClientOption) (*Client, error) {
	return &Client{hc: hc, opts: opts}, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string) (*http.Request, error) {
	path = strings.TrimLeft(path, "/")
	u := fmt.Sprintf("%s/%s", c.opts.Host, path)

	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", userAgent)

	expires := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{Name: sessionCookieName, Value: c.opts.SessionID, Expires: expires, HttpOnly: true, Secure: true}
	req.AddCookie(&cookie)

	return req, nil
}

func (c *Client) getGoQueryDoc(ctx context.Context, path string) (*goquery.Document, error) {
	req, err := c.newRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[mf] %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	if len(doc.Find("#header-container").Nodes) < 1 {
		return nil, fmt.Errorf("[mf] no login")
	}
	return doc, nil
}

func (c *Client) GetTransactionHistories(ctx context.Context) ([]*TransactionHistory, error) {
	doc, err := c.getGoQueryDoc(ctx, "/cf")
	if err != nil {
		return nil, err
	}

	var histories []*TransactionHistory
	doc.Find("#cf-detail-table tr.target-active").Each(func(i int, s *goquery.Selection) {
		date, err := time.Parse("2006/01/02", strings.Split(s.Find("td.date").AttrOr("data-table-sortable-value", ""), "-")[0])
		if err != nil {
			date = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		}
		content := strings.TrimSpace(s.Find("td.content").Text())
		amount, err := strconv.Atoi(strings.Replace(strings.TrimSpace(s.Find("td.amount").Text()), ",", "", -1))
		if err != nil {
			amount = 0
		}
		histories = append(histories, &TransactionHistory{
			Content: content,
			Amount:  amount,
			Date:    date,
		})
	})
	return histories, nil
}

// For debug
func (c *Client) GetHtml(ctx context.Context, spath string) (string, error) {
	doc, err := c.getGoQueryDoc(ctx, spath)
	if err != nil {
		return "", err
	}
	hs, err := doc.Html()
	if err != nil {
		return "", err
	}
	return hs, nil
}

func (c *Client) GetTotalAsset(ctx context.Context) (*TotalAsset, error) {
	doc, err := c.getGoQueryDoc(ctx, "/")
	if err != nil {
		return nil, err
	}

	txt := doc.Find(".total-assets .heading-radius-box").First().Text()
	re := regexp.MustCompile("[0-9]+")
	amount, err := strconv.Atoi(strings.Join(re.FindAllString(txt, -1), ""))
	if err != nil {
		return nil, err
	}
	return &TotalAsset{
		Amount: amount,
	}, nil
}
