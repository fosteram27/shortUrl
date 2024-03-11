package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fosteram27/shorturl/urls"
)

type EntriesHandler struct {
	store urls.Store
}

func NewEntriesHandler(s urls.Store) *EntriesHandler {
	return &EntriesHandler{
		store: s,
	}
}

func (h *EntriesHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var entry urls.UrlEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	log.Printf("creating short URL for %v", entry.UrlLong)

	// Add correct short URL
	hash := shortenURL(entry.UrlLong)
	entry.UrlShort = baseURL + hash

	if err := h.store.Add(hash, entry); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	// TODO: encode result as JSON
	w.Write([]byte(entry.UrlShort))

}

func (h *EntriesHandler) ListEntries(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()
	if err != nil {
		log.Println("first error listing entries", err)
		InternalServerErrorHandler(w, r)
		return
	}

	if err := json.NewEncoder(w).Encode(resources); err != nil {
		log.Println("second error listing entries", err)
		InternalServerErrorHandler(w, r)
		return
	}
}

func (h *EntriesHandler) GetEntry(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	log.Print("looking up URL for ", hash)

	// Retrieve recipe from the store
	entry, err := h.store.Get(hash)
	switch {
	case err == urls.ErrNotFound:
		NotFoundHandler(w, r)
		return
	case err != nil:
		InternalServerErrorHandler(w, r)
		return
	}

	// Convert the struct into JSON payload
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
}

func (h *EntriesHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	h.store.Remove(hash)

	w.WriteHeader(http.StatusOK)
}

func (h *EntriesHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	entry, err := h.store.Get(hash)
	switch {
	case err == urls.ErrNotFound:
		NotFoundHandler(w, r)
		return
	case err != nil:
		InternalServerErrorHandler(w, r)
		return
	}

	http.Redirect(w, r, entry.UrlLong, http.StatusFound)
}

// ERROR FUNCTIONS //
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func main() {
	// Create the store and recipe handler
	// store := urls.NewMemStore()
	store, err := urls.NewDBStore("urls.db")
	if err != nil {
		log.Fatalf("could not create DB: %v", err)
	}
	defer func() {
		log.Println("closing")
		store.Close()
	}()

	entriesHandler := NewEntriesHandler(store)

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is my home page"))
	}))
	//  add a new entry
	mux.Handle("GET /api", http.HandlerFunc(entriesHandler.ListEntries))
	mux.Handle("POST /api", http.HandlerFunc(entriesHandler.CreateEntry))
	mux.Handle("GET /api/{hash}", http.HandlerFunc(entriesHandler.GetEntry))
	mux.Handle("GET /{hash}", http.HandlerFunc(entriesHandler.Redirect))
	mux.Handle("DELETE /api/{hash}", http.HandlerFunc(entriesHandler.DeleteEntry))

	// mux.Handle("/0a137", &shortUrlHandler{})

	// Run the server
	log.Fatal(http.ListenAndServe(":8080", mux))
}

const (
	baseURL  = "bit.ly/"
	numChars = 5
)

func shortenURL(urlLong string) string {
	// TODO add timestamp before hashing to increase uniqueness of entries
	urlLongHash := md5.Sum([]byte(urlLong))
	//urlShort := baseURL + hex.EncodeToString(urlLongHash[:])[:numChars]
	return hex.EncodeToString(urlLongHash[:])[:numChars]
}
