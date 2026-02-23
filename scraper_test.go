package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScrapeSuccess(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from " + r.URL.Path))
	}))
	defer server.Close()
	
	// URLs a scrapear
	urls := []string{
		server.URL + "/page1",
		server.URL + "/page2",
		server.URL + "/page3",
	}
	
	// Crear scraper
	config := Config{
		MaxConcurrent:  2,
		RequestsPerSec: 10,
		Timeout:        5 * time.Second,
	}
	scraper := NewScraper(config)
	
	// Ejecutar scrape
	ctx := context.Background()
	results := scraper.Scrape(ctx, urls)
	
	// Verificar
	if len(results) != len(urls) {
		t.Fatalf("Expected %d results, got %d", len(urls), len(results))
	}
	
	for _, result := range results {
		if result.Error != nil {
			t.Errorf("URL %s failed: %v", result.URL, result.Error)
		}
		if result.Content == "" {
			t.Errorf("URL %s returned empty content", result.URL)
		}
	}
}

func TestScrapeWithTimeout(t *testing.T) {
	// Mock server que tarda 2 segundos
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("Slow response"))
	}))
	defer server.Close()
	
	urls := []string{server.URL}
	
	config := Config{
		MaxConcurrent:  1,
		RequestsPerSec: 10,
		Timeout:        500 * time.Millisecond,  // Timeout corto
	}
	scraper := NewScraper(config)
	
	ctx := context.Background()
	results := scraper.Scrape(ctx, urls)
	
	// Debe haber error de timeout
	if results[0].Error == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestScrapeWithCancellation(t *testing.T) {
	// Mock server lento
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("Response"))
	}))
	defer server.Close()
	
	urls := []string{
		server.URL + "/1",
		server.URL + "/2",
		server.URL + "/3",
	}
	
	config := Config{
		MaxConcurrent:  3,
		RequestsPerSec: 10,
		Timeout:        10 * time.Second,
	}
	scraper := NewScraper(config)
	
	// Context con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancelar después de 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	
	results := scraper.Scrape(ctx, urls)
	
	// No todos deberían completarse
	completed := 0
	for _, result := range results {
		if result.Error == nil {
			completed++
		}
	}
	
	if completed == len(urls) {
		t.Error("Expected some requests to be cancelled")
	}
}

func TestRateLimiting(t *testing.T) {
	requestTimes := []time.Time{}
	var mu sync.Mutex
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	
	// 10 URLs
	urls := make([]string, 10)
	for i := range urls {
		urls[i] = server.URL
	}
	
	config := Config{
		MaxConcurrent:  10,
		RequestsPerSec: 5,  // 5 requests/segundo
		Timeout:        5 * time.Second,
	}
	scraper := NewScraper(config)
	
	start := time.Now()
	ctx := context.Background()
	scraper.Scrape(ctx, urls)
	duration := time.Since(start)
	
	// Con rate de 5 req/s, 10 requests deberían tomar ~2 segundos
	// (burst inicial + rate limiting)
	if duration < 1*time.Second {
		t.Errorf("Too fast: %v (rate limiting not working)", duration)
	}
	
	t.Logf("Completed 10 requests in %v", duration)
}
