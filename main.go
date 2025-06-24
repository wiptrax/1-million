package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wiptrax/1-million-go/internal/db"
	"github.com/wiptrax/1-million-go/internal/sqlc"
)

type PostData struct {
	Timestamp int64 `json:"timestamp"`
}

func main() {
	// connect to db
	dbConn := db.ConnectToDB()

	qr := sqlc.New(dbConn)

	var wg sync.WaitGroup

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

		// Use channel or sync primitives if you want to track insert
		wg.Add(1)
		go func(data PostData) {
			defer wg.Done()
			if _, err := qr.CreateUser(context.TODO(), data.Timestamp); err != nil {
				log.Println("insert error:", err)
			}
		}(postdata)

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Waiting for all inserts to complete...")
		start :=  time.Now()
		wg.Wait()
		log.Println("All inserts done!\n", "time for all insert:",time.Since(start))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Done"))
	})

	log.Printf("listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
