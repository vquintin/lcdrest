package lcdrest

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

type RandomMessage struct {
	lock     sync.RWMutex
	messages map[string]string
	writer   io.Writer
	quit     chan int
}

func NewRandomMessage(w io.Writer, duration time.Duration) *RandomMessage {
	rm := &RandomMessage{
		messages: make(map[string]string),
		writer:   w,
		quit:     make(chan int),
	}
	go rm.writeRoutine(duration)
	return rm
}

func (rm *RandomMessage) Put(key string, message string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Put] Putting message '%v' for key '%v'.", message, key)
	defer log.Print("[lcdrest][RandomMessage][Put] Exit.")
	rm.lock.Lock()
	defer rm.lock.Unlock()
	old, exists := rm.messages[key]
	rm.messages[key] = message
	return old, !exists
}

func (rm *RandomMessage) Get(key string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Get] Getting message for key '%v'.", key)
	defer log.Print("[lcdrest][RandomMessage][Get] Exit.")
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	v, ok := rm.messages[key]
	return v, ok
}

func (rm *RandomMessage) Delete(key string) (string, bool) {
	log.Printf("[lcdrest][RandomMessage][Delete] Deleting message for key '%v'.", key)
	defer log.Print("[lcdrest][RandomMessage][Delete] Exit.")
	rm.lock.Lock()
	rm.lock.Unlock()
	old, ok := rm.messages[key]
	delete(rm.messages, key)
	return old, ok
}

func (rm *RandomMessage) Close() error {
	rm.quit <- 0
	return nil
}

func (rm *RandomMessage) writeRoutine(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			rm.writeRandomMessage()
		case <-rm.quit:
			return
		}
	}
}

func (rm *RandomMessage) writeRandomMessage() {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	n := int64(len(rm.messages))
	var r int64
	if n > 0 {
		r = rand.Int63n(n)
	}
	i := int64(0)
	var message string
	for _, v := range rm.messages {
		message = v
		if i == r {
			break
		}
		i++
	}
	rm.writer.Write([]byte(message))
}
