package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/edgar-macias-se/scraper"  // Ajusta el import según tu estructura
)

func main() {
	// URLs a scrapear
	urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://reddit.com",
	}
	
	// Configurar scraper
	config := scraper.Config{
		MaxConcurrent:  3,               // Máximo 3 requests simultáneos
		RequestsPerSec: 2,               // 2 requests por segundo
		Timeout:        10 * time.Second, // Timeout de 10s por request
	}
	
	s := scraper.NewScraper(config)
	
	// Context con timeout global
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Ejecutar scrape
	fmt.Println("Starting scraper...")
	start := time.Now()
	
	results := s.Scrape(ctx, urls)
	
	duration := time.Since(start)
	
	// Mostrar resultados
	fmt.Printf("\nCompleted in %v\n\n", duration)
	
	successful := 0
	failed := 0
	
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("❌ [%d] %s - Error: %v (took %v)\n",
				result.ID, result.URL, result.Error, result.Duration)
			failed++
		} else {
			fmt.Printf("✅ [%d] %s - %d bytes (took %v)\n",
				result.ID, result.URL, len(result.Content), result.Duration)
			successful++
		}
	}
	
	fmt.Printf("\nSummary: %d successful, %d failed\n", successful, failed)
}
