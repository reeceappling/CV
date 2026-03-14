package main

import "strings"

var interests = map[string]*Interest{}

type Interest struct {
	Name string
	*tracked
}

func NewInterest(name string) *Interest {
	out := &Interest{
		Name:    name,
		tracked: newTracked(),
	}
	interests[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}

func (pg *Interest) EntryType() string {
	return "Interest"
}

func (i *Interest) Link() string {
	return linkForInterest(i.Name)
}

func linkForInterest(name string) string {
	return linkFor(name, "cv", "interest", withoutSpaces(name))
}

var (
//	intMycology = NewInterest(string(
//
// smMycology)). // TODO: DEL?
// intCryptocurrency NewInterest(string(smCryptocurrency))
// intCryptography NewInterest(string(smCryptography))
)
