package main

import (
	"context"
	"log"
	"time"

	"github.com/Saumya40-codes/LogsGO/client/go/logclient"
)

var messages = []string{
	"Service timed out",
	"Connection refused",
	"Database unavailable",
	"Cache miss",
	"Authentication failed",
	"Context deadline exceeded",
}

var levels = []string{"info", "warn", "error"}
var services = []string{"ap-south1", "us-west1"}

// Timestamps for:
// - 1 Jan 2025 00:00:32 UTC
// - 30 May 2025 00:00:32 UTC
var fixedTimestamps = []int64{
	1746307232, // 1 Jan 2025
	1748553632, // 30 May 2025
}

func uploadLogWithDelay(lc *logclient.Client, opts *logclient.Opts) {
	ok := lc.UploadLog(opts)
	if !ok {
		log.Println("log upload failed")
	}
	log.Printf("Uploaded log: %+v", opts)
	time.Sleep(10 * time.Second)
}

func main() {
	ctx := context.Background()
	lc, err := logclient.NewLogClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// First, upload logs with fixed timestamps
	for i, ts := range fixedTimestamps {
		opts := &logclient.Opts{
			Level:     levels[i%len(levels)],
			Service:   services[i%len(services)],
			Message:   messages[i%len(messages)],
			TimeStamp: ts,
		}
		uploadLogWithDelay(lc, opts)
	}

	// Then, upload logs with time.Now()
	for i := 0; i < 5; i++ {
		for i := 0; i < 3; i++ {
			opts := &logclient.Opts{
				Level:     levels[i%len(levels)],
				Service:   services[i%len(services)],
				Message:   messages[i%len(messages)],
				TimeStamp: time.Now().Unix(),
			}

			uploadLogWithDelay(lc, opts)

			opts = &logclient.Opts{
				Level:     levels[(i+1)%len(levels)],
				Service:   services[(i+1)%len(services)],
				Message:   messages[(i+3)%len(messages)],
				TimeStamp: time.Now().Unix(),
			}
			uploadLogWithDelay(lc, opts)
		}

		time.Sleep(10 * time.Second)
	}
}
