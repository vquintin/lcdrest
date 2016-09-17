package randommessage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vquintin/lcdrest/messagestore"
)

type mockWriter struct {
	actual []string
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.actual = append(mw.actual, string(b))
	return len(b), nil
}

func TestWriteRandomMessage(t *testing.T) {
	messages := messagestore.NewSynchronizedMessageStore()
	expected := "some message"
	messages.Put("some key", expected)
	mw := mockWriter{}
	duration := 1 * time.Millisecond
	rm := NewRandomMessage(messages, &mw, duration)

	time.Sleep(duration)
	rm.Close()
	time.Sleep(duration)

	assert.Equal(t, 1, len(mw.actual))
	assert.Equal(t, expected, mw.actual[0])
}
