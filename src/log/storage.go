package log

type Storage struct {
	records map[Level][]record
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
	s.records[level] = append(s.records[level], record{message, caller})
}

func (s *Storage) Exists(level Level, message string, caller string) bool {
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
	s.records = map[Level][]record{}
}
