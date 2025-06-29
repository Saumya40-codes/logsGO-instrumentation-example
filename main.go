package main

import (
	"context"
	"log"
	"time"

	logapi "github.com/Saumya40-codes/LogsGO/api/grpc/pb"
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

var (
	levels   = []string{"info", "warn", "error"}
	services = []string{"ap-south1", "us-west1"}
)

var fixedTimestamps = []int64{
	1746307232, // goes to s3 (if configured)
	1748553632, // goes to local storage
}

func uploadLogWithDelay(ctx context.Context, lc *logclient.Client, opts *logapi.LogEntry) {
	err := lc.UploadLog(ctx, opts)
	if err != nil {
		log.Println("log upload failed")
	}
	log.Printf("Uploaded log: %+v", opts)
	time.Sleep(20 * time.Second)
}

func main() {
	time.Sleep(2 * time.Second)
	ctx := context.Background()
	lc, err := logclient.NewLogClient(ctx, "logsgo:50051")
	if err != nil {
		log.Fatal(err)
	}

	for i, ts := range fixedTimestamps {
		opts := &logapi.LogEntry{
			Level:     levels[i%len(levels)],
			Service:   services[i%len(services)],
			Message:   messages[i%len(messages)],
			Timestamp: ts,
		}
		uploadLogWithDelay(ctx, lc, opts)
	}

	for i := 0; i < 3; i++ { // goes in mem
		for i := 0; i < 3; i++ {
			opts := &logapi.LogEntry{
				Level:     levels[i%len(levels)],
				Service:   services[i%len(services)],
				Message:   messages[i%len(messages)],
				Timestamp: time.Now().Unix(),
			}

			uploadLogWithDelay(ctx, lc, opts)

			opts = &logapi.LogEntry{
				Level:     levels[(i+1)%len(levels)],
				Service:   services[(i+1)%len(services)],
				Message:   messages[(i+3)%len(messages)],
				Timestamp: time.Now().Unix(),
			}
			uploadLogWithDelay(ctx, lc, opts)
		}

		time.Sleep(20 * time.Second)
	}
}
