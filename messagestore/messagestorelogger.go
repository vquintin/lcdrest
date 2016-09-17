package messagestore

import "log"

type MessageStoreLogger struct {
	Delegate MessageStore
}

func (msl MessageStoreLogger) Put(key string, message string) (string, bool) {
	log.Printf("[messagestore][MessageStore][Put] Putting message '%v' for key '%v'.", message, key)
	return msl.Delegate.Put(key, message)
}

func (msl MessageStoreLogger) Get(key string) (string, bool) {
	log.Printf("[messagestore][MessageStore][Get] Getting message for key '%v'.", key)
	return msl.Delegate.Get(key)
}

func (msl MessageStoreLogger) GetAll() map[string]string {
	log.Printf("[messagestore][MessageStore][GetAll] Getting all messages.")
	return msl.Delegate.GetAll()
}

func (msl MessageStoreLogger) Delete(key string) (string, bool) {
	log.Printf("[messagestore][MessageStore][Delete] Deleting message for key '%v'.", key)
	return msl.Delegate.Delete(key)
}
