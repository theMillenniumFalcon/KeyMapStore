package main

import (
	"errors"
	"reflect"
	"strconv"
)

type Storage struct {
	Segments  map[rune]*Segment
	IdCounter int64
}

func (db *Storage) Init() {
	db.IdCounter = 0
	db.Segments = make(map[rune]*Segment)
}

func (db *Storage) Size() (size int) {
	size = 0
	for _, s := range db.Segments {
		size += s.Size
	}
	return size
}

func (db *Storage) Write(id string, data interface{}) error {
	if id == "" {
		return errors.New("passed ID is invalid")
	}

	if data == nil {
		return errors.New("passed data is nil")
	}

	var segID rune
	var seg *Segment
	var sub *Segment
	ok := false

	for i := 0; i < len(id); i++ {
		segID = rune(id[i])
		if seg == nil {
			if seg, ok = db.Segments[segID]; !ok {
				seg = &Segment{nil, 0, id[:i+1], make([]*Record, 0), make(map[rune]*Segment)}
				db.Segments[segID] = seg
			}
		} else {
			if sub, ok = seg.Sub[segID]; !ok {
				sub = &Segment{seg, 0, id[:i+1], make([]*Record, 0), make(map[rune]*Segment)}
				seg.Sub[segID] = sub
				seg = sub
			} else {
				seg = sub
			}
		}

	}

	if seg == nil {
		return errors.New("Segments were not evaluated properly for ID " + id)
	}

	db.IdCounter++
	seg.Records = append(seg.Records, &Record{db.IdCounter, id, data})
	// recursive incremental
	seg.Size++
	for seg.Root != nil {
		seg.Root.Size++
		seg = seg.Root
	}

	return nil
}

func (db *Storage) FindSegment(id string) (*Segment, error) {
	if id == "" {
		return nil, errors.New("passed ID is invalid")
	}

	var segmentID rune
	var seg *Segment
	ok := false

	for i := 0; i < len(id); i++ {
		segmentID = rune(id[i])
		if seg == nil {
			if seg, ok = db.Segments[segmentID]; !ok {
				return nil, errors.New("segment not found")
			}
		} else {
			if seg, ok = seg.Sub[segmentID]; !ok {
				return nil, errors.New("subSegment not found")
			}
		}

	}

	return seg, nil
}

// primaryFields must be type of string or int
func GetIdentifier(x interface{}, primaryFields []string) string {
	v := reflect.Indirect(reflect.ValueOf(x))
	t := v.Type()

	l := v.NumField()

	kinds := make([]reflect.Kind, l)
	values := make([]interface{}, l)
	keys := make([]reflect.StructField, l)

	identifier := ""

	for i := 0; i < l; i++ {
		kinds[i] = v.Field(i).Kind()
		values[i] = v.Field(i).Interface()
		keys[i] = t.Field(i)
	}

	temp := ""

	if primaryFields != nil {
		for _, field := range primaryFields {
			for i := 0; i < l; i++ {
				if field == keys[i].Name {
					if kinds[i] == reflect.String {
						temp = values[i].(string)
						if len(temp) > 0 {
							identifier += string(temp[0])
							break
						}
					} else if kinds[i] == reflect.Int {
						identifier += string(strconv.Itoa(values[i].(int))[0])
					}
				}
			}
		}

	} else {
		for i := 0; i < l; i++ {
			if kinds[i] == reflect.String {
				temp = values[i].(string)
				if len(temp) > 0 {
					identifier = string(temp[0])
					break
				}
			} else if kinds[i] == reflect.Int {
				identifier += string(strconv.Itoa(values[i].(int))[0])
			}
		}
	}

	return identifier
}

// primaryFields must be type of string or int
func GetIdentifierFromMap(m map[string]interface{}, primaryFields []string) string {
	identifier := ""
	temp := ""
	var k reflect.Kind

	if primaryFields != nil {
		for _, field := range primaryFields {
			v, ok := m[field]
			if ok {
				k = reflect.TypeOf(v).Kind()
				if k == reflect.String {
					temp = v.(string)
					if len(temp) > 0 {
						identifier += string(temp[0])
					}
				} else if k == reflect.Int {
					identifier += string(strconv.Itoa(v.(int))[0])
				}
			}
		}

	} else {
		for _, value := range m {
			if reflect.TypeOf(value).Kind() == reflect.String {
				temp = value.(string)
				if len(temp) > 0 {
					identifier = string(temp[0])
					break
				}
			}
		}
	}

	return identifier
}
