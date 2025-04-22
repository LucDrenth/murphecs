package log

import "sync"

type Storage struct {
	records map[Level][]record
	mutex   sync.RWMutex
}

type record struct {
	message string
	caller  string // from where in the code the log was called
}

func NewStorage() Storage {
	return Storage{
		records: map[Level][]record{},
	}
}

func (s *Storage) Insert(level Level, message string, caller string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.records[level] = append(s.records[level], record{message, caller})
}

func (s *Storage) Exists(level Level, message string, caller string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	records, exists := s.records[level]
	if !exists {
		return false
	}

	for i := range records {
		if records[i].caller == caller && records[i].message == message {
			return true
		}
	}

	return false
}

func (s *Storage) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.records = map[Level][]record{}
}
