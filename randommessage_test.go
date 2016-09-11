package lcdrest

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPutAndGetRandomMessage(t *testing.T) {
	var buf bytes.Buffer
	rm := NewRandomMessage(&buf, 15*time.Second)
	key := "some key"
	expected := "some message"

	_, created := rm.Put(key, expected)
	actual, ok := rm.Get(key)

	assert.True(t, created)
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestPutAndDeleteRandomMessage(t *testing.T) {
	var buf bytes.Buffer
	rm := NewRandomMessage(&buf, 15*time.Second)
	key := "some key"
	expected := "some message"

	_, created := rm.Put(key, expected)
	deletedValue, found := rm.Delete(key)

	assert.True(t, created)
	assert.True(t, found)
	assert.Equal(t, expected, deletedValue)
}

func TestDoublePutRandomMessage(t *testing.T) {
	var buf bytes.Buffer
	rm := NewRandomMessage(&buf, 15*time.Second)
	key := "some key"
	expectedOld := "some message"
	expectedNew := "new message"

	_, createdOld := rm.Put(key, expectedOld)
	actualOld, createdNew := rm.Put(key, expectedNew)
	actualNew, _ := rm.Get(key)

	assert.True(t, createdOld)
	assert.Equal(t, expectedOld, actualOld)
	assert.False(t, createdNew)
	assert.Equal(t, expectedNew, actualNew)
}

type mockWriter struct {
	lock   sync.Mutex
	rm     *RandomMessage
	actual string
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.actual = string(b)
	return len(b), nil
}

func TestWriteRandomMessage(t *testing.T) {
	mw := mockWriter{}
	duration := 1 * time.Millisecond
	rm := NewRandomMessage(&mw, duration)
	key := "some key"
	expected := "some message"

	rm.Put(key, expected)

	time.Sleep(2 * duration)
	mw.lock.Lock()
	assert.Equal(t, expected, mw.actual)
}
