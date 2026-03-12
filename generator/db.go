package main

import "fmt"

var dbs = map[string]*DbPage{}

type DbPage struct {
	Name string
	*tracked
}

func (pg *DbPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "db", withoutSpaces(pg.Name))
}

func NewDb(name string) *DbPage {
	out := &DbPage{
		Name:    name,
		tracked: newTracked(),
	}
	dbs[name] = out
	return out
}
