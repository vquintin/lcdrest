package messagestore

import (
	"log"
	"sync"
)

type MessageStore interface {
	Put(key string, message string) (string, bool)
	Get(key string) (string, bool)
	GetAll() map[string]string
	Delete(key string) (string, bool)
}

type SynchronizedMessageStore struct {
	lock     sync.RWMutex
	messages map[string]string
}

func (sms *SynchronizedMessageStore) Put(key string, message string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Put] Putting message '%v' for key '%v'.", message, key)
	defer log.Print("[lcdrest][RandomMessage][Put] Exit.")
	sms.lock.Lock()
	defer sms.lock.Unlock()
	old, exists := sms.messages[key]
	sms.messages[key] = message
	return old, !exists
}

func (sms *SynchronizedMessageStore) Get(key string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Get] Getting message for key '%v'.", key)
	defer log.Print("[lcdrest][RandomMessage][Get] Exit.")
	sms.lock.RLock()
	defer sms.lock.RUnlock()
	v, ok := sms.messages[key]
	return v, ok
}

func (sms *SynchronizedMessageStore) GetAll() map[string]string {
	sms.lock.RLock()
	defer sms.lock.RUnlock()
	copy := make(map[string]string, len(sms.messages))
	for k, v := range sms.messages {
		copy[k] = v
	}
	return copy
}

func (sms *SynchronizedMessageStore) Delete(key string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Delete] Deleting message for key '%v'.", key)
	defer log.Print("[lcdrest][RandomMessage][Delete] Exit.")
	sms.lock.Lock()
	sms.lock.Unlock()
	old, ok := sms.messages[key]
	delete(sms.messages, key)
	return old, ok
}

func NewSynchronizedMessageStore() *SynchronizedMessageStore {
	rm := &SynchronizedMessageStore{
		messages: make(map[string]string),
	}
	return rm
}
