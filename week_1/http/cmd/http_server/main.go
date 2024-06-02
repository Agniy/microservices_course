package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/olezhek28/microservices_course/week_1/http/pkg/membd"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/note/{noteId}"
)

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &membd.NoteInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to decode note data", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	noteObj := &membd.Note{
		ID:        rand.Int63(),
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(noteObj); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}

	membd.Notes.M.Lock()
	defer membd.Notes.M.Unlock()

	membd.Notes.Elems[noteObj.ID] = noteObj
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteId")
	id, err := parseNoteID(noteID)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	membd.Notes.M.RLock()
	defer membd.Notes.M.RUnlock()

	note, ok := membd.Notes.Elems[id]

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", note.ID, note.Info.Title, note.Info, note.CreatedAt, note.UpdatedAt)

	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
}

func parseNoteID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func main() {
	r := chi.NewRouter()

	r.Get(getPostfix, getNoteHandler)
	r.Post(createPostfix, createNoteHandler)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
