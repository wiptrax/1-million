package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wiptrax/1-million-go/internal/db"
)

type PostData struct {
	Timestamp int64 `json:"timestamp"`
}

const (
	workerCount   = 10    // 10 workers with batch insert = faster than 50 single-row workers
	jobQueueSize  = 10000 // buffered job queue
	batchSize     = 100   // number of inserts per batch
)

func main() {
	dbConn := db.ConnectToDB()
	jobQueue := make(chan PostData, jobQueueSize)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go worker(i, dbConn, jobQueue, &wg)
	}

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "not a post method", http.StatusMethodNotAllowed)
			return
		}

		var postdata PostData
		err := json.NewDecoder(r.Body).Decode(&postdata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wg.Add(1)
		jobQueue <- postdata
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Waiting for all inserts to complete...")
		start := time.Now()
		wg.Wait()
		log.Println("All inserts done!")
		log.Println("time for all insert:", time.Since(start))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Done"))
	})

	log.Printf("listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Worker batches and inserts
func worker(id int, dbConn *pgxpool.Pool, jobs <-chan PostData, wg *sync.WaitGroup) {
	batch := make([]PostData, 0, batchSize)

	flush := func(b []PostData) {
		if len(b) == 0 {
			return
		}
		err := batchInsert(dbConn, b)
		if err != nil {
			log.Printf("Worker %d batch insert error: %v\n", id, err)
		}
		for range b {
			wg.Done()
		}
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case job := <-jobs:
			batch = append(batch, job)
			if len(batch) >= batchSize {
				flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flush(batch)
				batch = batch[:0]
			}
		}
	}
}

// Batch insert using pgx
func batchInsert(db *pgxpool.Pool, data []PostData) error {
	if len(data) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data))

	for i, row := range data {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i+1))
		args = append(args, row.Timestamp)
	}

	query := fmt.Sprintf("INSERT INTO users (insert_time_milli) VALUES %s", strings.Join(valueStrings, ","))
	_, err := db.Exec(context.Background(), query, args...)
	return err
}
