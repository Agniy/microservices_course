package membd

import (
	"sync"
	"time"
)

var Notes *SyncMap

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SyncMap struct {
	Elems map[int64]*Note
	M     sync.RWMutex
}

func init() {
	Notes = &SyncMap{
		Elems: make(map[int64]*Note),
	}
}
