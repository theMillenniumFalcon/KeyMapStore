package main

type Segment struct {
	Root    *Segment `json:"-"`
	Size    int
	Path    string
	Records []*Record
	Sub     map[rune]*Segment
}
