package memory

import (
	"TestOzon/internal/handler/graph/model"
	"github.com/google/uuid"
	"sync"
)

type StoreMemory struct {
	mu       sync.RWMutex
	Posts    map[uuid.UUID]*model.Post
	Comments map[uuid.UUID][]*model.Comment
}

func InitMemory() *StoreMemory {
	return &StoreMemory{
		Posts:    make(map[uuid.UUID]*model.Post),
		Comments: make(map[uuid.UUID][]*model.Comment),
	}
}
