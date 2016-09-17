package randommessage

import (
	"math/rand"

	"github.com/vquintin/lcdrest/messagestore"
)

type RandomMessageGetter struct {
	ms messagestore.MessageStore
}

func (rmg RandomMessageGetter) RandomMessage() (string, bool) {
	messages := rmg.ms.GetAll()
	n := int64(len(messages))
	var r int64
	if n > 0 {
		r = rand.Int63n(n)
		i := int64(0)
		for _, v := range messages {
			if i == r {
				return v, true
			}
			i++
		}

	}
	return "", false
}
