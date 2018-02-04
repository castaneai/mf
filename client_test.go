package mf

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

func newTestClient() (*Client, error) {
	hc := &http.Client{}
	opts := &ClientOption{Host: "https://moneyforward.com", SessionID: os.Getenv("MF_SESSION")}
	c, err := NewClient(hc, opts)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestClient_GetTransactionHistories(t *testing.T) {
	c, err := newTestClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	th, err := c.GetTransactionHistories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, trans := range th {
		t.Logf("%+v", trans)
	}
}

func TestClient_GetTotalAsset(t *testing.T) {
	c, err := newTestClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	ta, err := c.GetTotalAsset(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", ta)
}

func TestClient_NoLogin(t *testing.T) {
	hc := &http.Client{}
	opts := &ClientOption{Host: "https://moneyforward.com", SessionID: "invalid session"}
	c, err := NewClient(hc, opts)
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if _, err := c.GetTotalAsset(ctx); err != nil {
		t.Logf("%+v", err)
	}
}
