package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Job representa una URL a scrapear
type Job struct {
	URL string
	ID  int  // Para tracking
}

// Result representa el resultado del scrape
type Result struct {
	URL      string
	Content  string
	Error    error
	Duration time.Duration
	ID       int
}

// Config configura el scraper
type Config struct {
	MaxConcurrent int           // Máximo workers simultáneos
	RequestsPerSec int          // Rate limit (requests/segundo)
	Timeout       time.Duration // Timeout por request
}

// Scraper es el web scraper concurrente
type Scraper struct {
	config      Config
	rateLimiter *RateLimiter
	semaphore   *Semaphore
	client      *http.Client
}

// NewScraper crea un nuevo scraper
func NewScraper(config Config) *Scraper {
	return &Scraper{
		config:      config,
		rateLimiter: NewRateLimiter(
			time.Second/time.Duration(config.RequestsPerSec),
			config.RequestsPerSec*2,  // Burst = 2x rate
		),
		semaphore: NewSemaphore(config.MaxConcurrent),
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Scrape ejecuta el scraping de múltiples URLs
func (s *Scraper) Scrape(ctx context.Context, urls []string) []Result {
	jobs := make(chan Job, len(urls))
	results := make(chan Result, len(urls))

	var wg sync.WaitGroup

	numWorkers := s.config.MaxConcurrent
	for i := 0; i < numWorkers; i++{
		wg.Add(1)
		go func(workerID int){
			defer wg.Done()
			s.worker(ctx, workerID, jobs, results)
		}(i)
	}

	for i, url := range urls {
		jobs <- Job{URL: url, ID:i}
	}
	close(jobs)

	go func(){
			wg.Wait()
			close(results)
	}()

	var allResults []Result
	for result := range results{
		allResults = append(allResults, result)
	}

	return allResults
}

// worker procesa jobs individuales
func (s *Scraper) worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result) {
		
	for job := range jobs {
		select{
			case <-ctx.Done():
				return
			default:
		}
		
		s.semaphore.Acquire()

		s.rateLimiter.Wait()

		start := time.Now()
		content, err := s.fetch(ctx, job.URL)

		results <- Result{
							URL: job.URL,
							Content: content,
							Error: err,
							Duration: time.Since(start),
							ID: job.ID,
					}
		}

		s.semaphore.Release()
}

// fetch descarga una URL
func (s *Scraper) fetch(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    resp, err := s.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    
    return string(body), nil
}
