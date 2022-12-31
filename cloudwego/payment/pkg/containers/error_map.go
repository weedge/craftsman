// just for concurency err map

package containers

import (
	"fmt"
	"sync"
)

type ErrMap struct {
	mapErr map[string]error
	mu     sync.RWMutex
}

func NewErrMap() ErrMap {
	return ErrMap{
		mapErr: map[string]error{},
	}
}

func (m *ErrMap) Add(key string, err error) {
	m.mu.Lock()
	m.mapErr[key] = err
	m.mu.Unlock()
}

func (m *ErrMap) Get(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err, ok := m.mapErr[key]; ok {
		return err
	}
	return nil
}

func (m *ErrMap) Remove(key string) {
	m.mu.Lock()
	delete(m.mapErr, key)
	m.mu.Unlock()
}

func (m *ErrMap) Len() int {
	return len(m.mapErr)
}

func (m *ErrMap) String() string {
	str := "|"
	for k, err := range m.mapErr {
		if err != nil {
			str += fmt.Sprintf(" %s : %s |", k, err.Error())
		}
	}
	return str
}
