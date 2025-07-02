package urls_storage

import (
	"errors"
	"sync"
)

type Repository map[string]string

type SafeRepository struct {
	mu   sync.RWMutex
	data Repository
}

func NewSafeRepository() *SafeRepository {
	return &SafeRepository{
		data: make(Repository),
	}
}

func (sm *SafeRepository) Set(key string, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeRepository) Get(key string) (string, bool) {
	val, ok := sm.data[key]
	return val, ok
}

var urlsStorage = NewSafeRepository()

func mGetFullURL(shortURL string) (string, error) {
	fullURL, exists := urlsStorage.Get(shortURL)
	if !exists {
		return "", errors.New("URL doesnt exists")
	}
	return fullURL, nil
}

func mSetFullURL(fullURL string) (string, error) {
	shortURL := makeShortURL(fullURL)
	urlsStorage.Set(shortURL, fullURL)
	return shortURL, nil
}
