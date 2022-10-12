package inmemory

import (
	"main/internal/log"
	"main/internal/model"
	"sync"
)

// Cache struct holds cached data.
type Cache struct {
	Storage map[string]model.RecordModel
	Mutex   *sync.Mutex
}

// NewCache return new instance of Cache.
func NewCache() *Cache {
	return &Cache{
		Storage: make(map[string]model.RecordModel),
		Mutex:   &sync.Mutex{},
	}
}

// Insert loads model to cache.
func (c Cache) Insert(record model.RecordModel) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, ok := c.Storage[record.OrderUID]; ok {
		log.Logger.Traceln("Insert: message is not unique")
		return
	}

	c.Storage[record.OrderUID] = record

	return
}
