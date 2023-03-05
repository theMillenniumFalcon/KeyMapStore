package main

type Storage struct {
	Segments  map[rune]*Segment
	IDCounter int64
}
