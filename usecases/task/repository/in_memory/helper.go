package inmemory

import "sync"

func getLen(m *sync.Map) int {
	len := 0
	m.Range(func(k, v interface{}) bool {
		len++
		return true
	})
	return len
}
