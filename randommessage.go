package lcdrest

import (
	"io"
	"log"
	"math/rand"
	"time"
)

type pair struct {
	key     string
	message string
}

type randomMessage struct {
	writer     io.Writer
	add        chan pair
	delete     chan string
	getRequest chan int
	getAnswer  chan map[string]string
	quit       chan int
}

func NewRandomMessage(w io.Writer) randomMessage {
	rm := randomMessage{writer: w}
	go monitor(rm)
	return rm
}

func (rm randomMessage) Add(key string, message string) {
	log.Printf("[lcdrest][randomMessage][Add] Adding message '%v' for key '%v'.", message, key)
	rm.add <- pair{
		key:     key,
		message: message,
	}
	log.Print("[lcdrest][randomMessage][Add] Exit.")
}

func (rm randomMessage) Delete(key string) {
	rm.delete <- key
}

func (rm randomMessage) GetAll() map[string]string {
	rm.getRequest <- 0
	return <-rm.getAnswer
}

func (rm randomMessage) Close() error {
	rm.quit <- 0
	return nil
}

func monitor(rm randomMessage) {
	messages := make(map[string]string)
	ticker := time.NewTicker(15 * time.Second)
	select {
	case m := <-rm.add:
		messages[m.key] = m.message
	case k := <-rm.delete:
		delete(messages, k)
	case <-rm.getRequest:
		messagesCopy := copyMessages(messages)
		rm.getAnswer <- messagesCopy
	case <-ticker.C:
		writeRandomMessage(rm.writer, messages)
	case <-rm.quit:
		return
	}
}

func copyMessages(messages map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range messages {
		newMap[k] = v
	}
	return newMap
}

func writeRandomMessage(writer io.Writer, messages map[string]string) {
	n := int64(len(messages))
	var r int64 = 0
	if n > 0 {
		r = rand.Int63n(n)
	}
	i := int64(0)
	var message string
	for _, v := range messages {
		message = v
		if i == r {
			break
		}
		i++
	}
	writer.Write([]byte(message))
}
