package main

import (
	"context"
	"flag"
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
	levels      = []string{"info", "warn", "error"}
	services    = []string{"ap-south1", "us-west1"}
	enableQueue = false
)

var fixedTimestamps = []int64{
	1746307232, // goes to s3 (if configured)
	1748553632, // goes to local storage
}

func init() {
	flag.BoolVar(&enableQueue, "enable-queue", false, "Use if logs should be uploaded to rabbitmq server")
	flag.Parse()
}

func uploadLogWithDelay(ctx context.Context, lc *logclient.Client, opts *logapi.LogEntry) {
	err := lc.UploadLog(ctx, opts)
	if err != nil {
		log.Println("log upload failed")
	}
	log.Printf("Uploaded log: %+v", opts)
	time.Sleep(20 * time.Second)
}

func uploadBatchToQueue(ctx context.Context, lc *logclient.Client, entries []*logapi.LogEntry) {
	batch := &logapi.LogBatch{
		Entries: entries,
	}
	err := lc.UploadLogsToQueue(ctx, batch)
	if err != nil {
		log.Println("queue log upload failed:", err)
	}
	log.Printf("Uploaded batch to queue: %d logs\n", len(entries))
	time.Sleep(20 * time.Second)
}

func main() {
	time.Sleep(2 * time.Second)
	ctx := context.Background()

	var lc *logclient.Client
	var err error

	if enableQueue {
		lc, err = logclient.NewLogClientWithQueue(ctx, "logsgo:50051", &logclient.QueueOpts{
			Url:       "amqp://guest:guest@rabbitmq:5672/",
			QueueName: "logs",
		}, true, nil)
	} else {
		lc, err = logclient.NewLogClient(ctx, "logsgo:50051")
	}
	if err != nil {
		log.Fatal(err)
	}

	if enableQueue {
		var batch []*logapi.LogEntry

		// Add fixed timestamp logs
		for i, ts := range fixedTimestamps {
			batch = append(batch, &logapi.LogEntry{
				Level:     levels[i%len(levels)],
				Service:   services[i%len(services)],
				Message:   messages[i%len(messages)],
				Timestamp: ts,
			})
		}

		// Add generated logs
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				batch = append(batch, &logapi.LogEntry{
					Level:     levels[j%len(levels)],
					Service:   services[j%len(services)],
					Message:   messages[j%len(messages)],
					Timestamp: time.Now().Unix(),
				})

				batch = append(batch, &logapi.LogEntry{
					Level:     levels[(j+1)%len(levels)],
					Service:   services[(j+1)%len(services)],
					Message:   messages[(j+3)%len(messages)],
					Timestamp: time.Now().Unix(),
				})
			}
			uploadBatchToQueue(ctx, lc, batch)
			batch = nil // Clear batch for next round
		}
	} else {
		// No queue, use gRPC upload
		for i, ts := range fixedTimestamps {
			opts := &logapi.LogEntry{
				Level:     levels[i%len(levels)],
				Service:   services[i%len(services)],
				Message:   messages[i%len(messages)],
				Timestamp: ts,
			}
			uploadLogWithDelay(ctx, lc, opts)
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				opts := &logapi.LogEntry{
					Level:     levels[j%len(levels)],
					Service:   services[j%len(services)],
					Message:   messages[j%len(messages)],
					Timestamp: time.Now().Unix(),
				}
				uploadLogWithDelay(ctx, lc, opts)

				opts = &logapi.LogEntry{
					Level:     levels[(j+1)%len(levels)],
					Service:   services[(j+1)%len(services)],
					Message:   messages[(j+3)%len(messages)],
					Timestamp: time.Now().Unix(),
				}
				uploadLogWithDelay(ctx, lc, opts)
			}
			time.Sleep(20 * time.Second)
		}
	}
}
