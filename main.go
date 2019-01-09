package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Run http server
func run() error {
	r := mux.NewRouter()
	r.HandleFunc("/", handleGetBlockchain).Methods("GET")
	r.HandleFunc("/", handleWriteBlock).Methods("POST")
	serverAddr := os.Getenv("SERVER_ADDR")
	log.Println("Listening on", serverAddr)
	s := &http.Server{
		Addr:           serverAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	return err
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var p Payload

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&p); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], p.Data)
	mutex.Unlock()

	if isBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
		Blockchain = append(Blockchain, newBlock)
	}

	respondWithJSON(w, r, http.StatusAccepted, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

var envfile = flag.String("env-file", ".env", "File containing environment variables.")
var mutex = &sync.Mutex{}

func main() {
	flag.Parse()
	if err := godotenv.Load(*envfile); err != nil {
		log.Fatal(err)
	}

	go func() {
		mutex.Lock()
		Blockchain = append(Blockchain, generateBlock(Block{}, "first block"))
		mutex.Unlock()
	}()

	log.Fatal(run())
}
