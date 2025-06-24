package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

	user, err := qr.CreateUser(context.TODO(), postdata.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log.Println("Inserted user:", user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	})	

	log.Printf("listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}