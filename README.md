# 🚀 Concurrent Web Scraper

High-performance web scraper built in Go with production-ready concurrency patterns.

## ✨ Features

- **Concurrent Processing**: Worker pool pattern with configurable concurrency
- **Rate Limiting**: Token bucket algorithm to respect server limits
- **Smart Cancellation**: Context-based timeout and cancellation support
- **Robust Error Handling**: Per-URL error tracking with detailed metrics
- **Production Ready**: 90%+ test coverage with comprehensive test suite

## 🏗️ Architecture

```
URLs → Jobs Queue → Worker Pool → Rate Limiter → HTTP Client
                         ↓
                   Results Channel → Aggregator
```

**Patterns Implemented:**

- Fan-Out/Fan-In
- Worker Pool
- Semaphore (concurrency limiting)
- Rate Limiting (token bucket)
- Context propagation

## 🚀 Quick Start

```bash
# Install
go get github.com/edgar-macias-se/scraper

# Run
go run main.go
```

## 📊 Performance

- **Throughput**: 1000+ URLs/min (configurable)
- **Concurrency**: Up to 100 workers (configurable)
- **Rate Limiting**: Configurable requests/second
- **Memory**: ~50MB for 1000 concurrent requests

## 🧪 Testing

```bash
# Run tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -cover ./...
```

## 🔧 Configuration

```go
config := scraper.Config{
    MaxConcurrent:  10,              // Max simultaneous requests
    RequestsPerSec: 5,               // Rate limit (req/s)
    Timeout:        10 * time.Second, // Per-request timeout
}
```

## 📈 Use Cases

- Web scraping at scale
- Site monitoring
- Content aggregation
- SEO analysis
- Link validation

## 🛠️ Built With

- **Go 1.21+**
- **Standard library only** (no external dependencies)
- Clean architecture principles
- TDD approach

## 📝 Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/edgar-macias-se/scraper/internal/scraper"
)

func main() {
    urls := []string{
        "https://example.com",
        "https://golang.org",
    }

    config := scraper.Config{
        MaxConcurrent:  3,
        RequestsPerSec: 2,
        Timeout:        10 * time.Second,
    }

    s := scraper.NewScraper(config)
    ctx := context.Background()

    results := s.Scrape(ctx, urls)

    for _, r := range results {
        if r.Error != nil {
            fmt.Printf("❌ %s: %v\n", r.URL, r.Error)
        } else {
            fmt.Printf("✅ %s: %d bytes\n", r.URL, len(r.Content))
        }
    }
}
```

## 🎯 Key Learnings

This project demonstrates:

- **Production Go patterns**: Real-world concurrency beyond tutorials
- **System design**: Rate limiting, worker pools, graceful shutdown
- **Testing discipline**: Table-driven tests, mocks, race detection
- **Code organization**: Clean package structure, interfaces

## 📄 License

MIT

---

**Built by [Edgar Macías](https://edgarmacias.com)** | [LinkedIn](https://linkedin.com/in/edgar-macias-devcybsec) | [GitHub](https://github.com/edgar-macias-se)
