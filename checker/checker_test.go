package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := CheckURL(server.URL, 2*time.Second)

	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", result.StatusCode)
	}
}

func TestCheckURL_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := CheckURL(server.URL, 1*time.Second)

	if result.Err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestCheckAll_ReturnsAllResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL, server.URL, server.URL, server.URL}
	results := CheckAll(urls, 2, 2*time.Second)

	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
}