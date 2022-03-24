package topics

import (
	"sync"
	"time"
)

const DEFAULT_TTL = 20000

type Topics struct {
	rw   *sync.RWMutex
	data map[string]time.Time
	ttl  time.Duration
}

func (t *Topics) GetValids() []string {
	valids := []string{}
	timeNow := time.Now()
	t.rw.RLock()
	for topic, timestamp := range t.data {
		if timestamp.After(timeNow) {
			valids = append(valids, topic)
		}
	}
	t.rw.RUnlock()
	return valids
}

func (t *Topics) UpdateTopic(topic string) {
	t.rw.Lock()
	timeNow := time.Now()
	t.data[topic] = timeNow.Add(time.Millisecond * t.ttl)
	t.rw.Unlock()
}

var once sync.Once
var instance *Topics

func Get() *Topics {
	once.Do(func() {
		instance = &Topics{
			rw:   &sync.RWMutex{},
			data: make(map[string]time.Time),
			ttl:  DEFAULT_TTL,
		}
	})
	return instance
}
