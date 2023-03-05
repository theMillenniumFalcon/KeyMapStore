package main

import "reflect"

type Segment struct {
	Root    *Segment `json:"-"`
	Size    int
	Path    string
	Records []*Record
	Sub     map[rune]*Segment
}

func (seg *Segment) FindFirstRecord(pattern map[string]interface{}) *Record {
	var rV reflect.Value
	var temp reflect.Value

	hits := 0
	desiredAmount := len(pattern)

	for _, r := range seg.Records {
		hits = 0

		rV = reflect.ValueOf(r.Data)

		for k, v := range pattern {
			temp = rV.FieldByName(k)
			if temp.IsValid() {
				if reflect.DeepEqual(v, temp.Interface()) {
					hits++
				}
			}
		}

		if hits >= desiredAmount {
			return r
		}
	}

	return nil
}

func (seg *Segment) FindRecords(pattern map[string]interface{}) []*Record {
	var rV reflect.Value
	var temp reflect.Value

	hits := 0
	desiredAmount := len(pattern)

	records := make([]*Record, 0)

	for _, r := range seg.Records {
		hits = 0

		rV = reflect.ValueOf(r.Data)

		for k, v := range pattern {
			temp = rV.FieldByName(k)
			if temp.IsValid() {
				if reflect.DeepEqual(v, temp.Interface()) {
					hits++
				}
			}
		}

		if hits >= desiredAmount {
			records = append(records, r)
		}
	}

	return records
}
