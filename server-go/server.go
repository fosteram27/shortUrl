package main

import (
	"api-go/urls"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
)

// REGEX for validating inputs
var (
	EntryRe = regexp.MustCompile(`^/entries/*$`)
	//EntryReWithID = regexp.MustCompile(`^/entries/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
	EntryReWithID = regexp.MustCompile(`^/entries/([a-z0-9]+)$`)
)

// storage, replace with database in future versions
type entryStore interface {
	Add(name string, entry urls.UrlEntry) error
	Get(name string) (urls.UrlEntry, error)
	Update(name string, entry urls.UrlEntry) error
	List() (map[string]urls.UrlEntry, error)
	Remove(name string) error
}

type entryDb interface {
	Add(name string, entry urls.UrlEntry) error
	Get(name string) (urls.UrlEntry, error)
	Update(name string, entry urls.UrlEntry) error
	List() (map[string]urls.UrlEntry, error)
	Remove(name string) error
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

// type shortUrlHandler struct{}

// func (h *shortUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// look up longUrl from shortUrl
// 	http.Redirect(w, r, "https://google.com", http.StatusSeeOther)

// }

type EntriesHandler struct {
	store entryStore
}

func NewEntriesHandler(s entryStore) *EntriesHandler {
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

	//Add correct short URL
	entry.UrlShort = shortenUrl(entry.UrlLong)

	//TODO: Modify later, need to figure out how to ID entries. longUrl?
	resourceID := entry.UrlShort

	err = h.store.Add(resourceID, entry)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.Write([]byte(entry.UrlShort))
	// w.WriteHeader(http.StatusOK)
}

func (h *EntriesHandler) ListEntries(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *EntriesHandler) GetEntry(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID/slug using a regex
	matches := EntryReWithID.FindStringSubmatch(r.URL.Path)
	// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve recipe from the store
	entry, err := h.store.Get(matches[1])
	if err != nil {
		// Special case of NotFound Error
		if err == urls.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		// Every other error
		InternalServerErrorHandler(w, r)
		return
	}

	// Convert the struct into JSON payload
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Write the results
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h *EntriesHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	matches := EntryReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Recipe object that will be populated from JSON payload
	var entry urls.UrlEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(matches[1], entry); err != nil {
		if err == urls.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *EntriesHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	matches := EntryReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EntriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == http.MethodPost && EntryRe.MatchString(r.URL.Path):
		h.CreateEntry(w, r)
		return
	case r.Method == http.MethodGet && EntryRe.MatchString(r.URL.Path):
		h.ListEntries(w, r)
		return
	case r.Method == http.MethodGet && EntryReWithID.MatchString(r.URL.Path):
		h.GetEntry(w, r)
		return
	case r.Method == http.MethodPut && EntryReWithID.MatchString(r.URL.Path):
		h.UpdateEntry(w, r)
		return
	case r.Method == http.MethodDelete && EntryReWithID.MatchString(r.URL.Path):
		h.DeleteEntry(w, r)
		return
	default:
		return
	}
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

// MAIN //
func main() {

	//Create the store and recipe handler
	store := urls.NewMemStore()
	entriesHandler := NewEntriesHandler(store)

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/entries", entriesHandler)
	mux.Handle("/entries/", entriesHandler)
	// mux.Handle("/0a137", &shortUrlHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)

}

func shortenUrl(urlLong string) string {

	url_base := "bit.ly/"
	numChars := 5

	//TODO add timestamp before hashing to increase uniqueness of entries
	urlLongHash := md5.Sum([]byte(urlLong))
	urlShort := url_base + hex.EncodeToString(urlLongHash[:])[:numChars]
	return urlShort
}
