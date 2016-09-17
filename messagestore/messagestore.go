package messagestore

type MessageStore interface {
	Put(key string, message string) (string, bool)
	Get(key string) (string, bool)
	GetAll() map[string]string
	Delete(key string) (string, bool)
}
