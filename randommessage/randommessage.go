package randommessage

import (
	"io"
	"time"

	"github.com/vquintin/lcdrest/messagestore"
)

type RandomMessage struct {
	rmg    RandomMessageGetter
	writer io.Writer
	quit   chan int
}

func (rm RandomMessage) Close() error {
	rm.quit <- 0
	return nil
}

func (rm RandomMessage) writeRoutine(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			message, ok := rm.rmg.RandomMessage()
			if ok {
				rm.writer.Write([]byte(message))
			}
		case <-rm.quit:
			return
		}
	}
}

func NewRandomMessage(messages messagestore.MessageStore, w io.Writer, duration time.Duration) RandomMessage {
	rm := RandomMessage{
		rmg:    RandomMessageGetter{messages},
		writer: w,
		quit:   make(chan int),
	}
	go rm.writeRoutine(duration)
	return rm
}
