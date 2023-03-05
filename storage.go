package main

type Storage struct {
	Segments  map[rune]*Segment
	IDCounter int64
}

func (db *Storage) Init() {
	db.IDCounter = 0
	db.Segments = make(map[rune]*Segment)
}

func (db *Storage) Size() (size int) {
	size = 0
	for _, s := range db.Segments {
		size += s.Size
	}

	return size
}
