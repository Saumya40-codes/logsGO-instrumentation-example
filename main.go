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

var fixedTimestamps = []int64{
	1746307232, // goes to s3 (if configured)
	1748553632, // goes to local storage
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

	for i, ts := range fixedTimestamps {
		opts := &logclient.Opts{
			Level:     levels[i%len(levels)],
			Service:   services[i%len(services)],
			Message:   messages[i%len(messages)],
			TimeStamp: ts,
		}
		uploadLogWithDelay(lc, opts)
	}

	for i := 0; i < 3; i++ { // goes in mem
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
