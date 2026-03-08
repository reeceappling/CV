package main

import "fmt"

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
	return out
}

func (i *Interest) Link() string {
	return linkForInterest(i.Name)
}

func linkForInterest(name string) string {
	return fmt.Sprintf("[%s](interest/%s)", name, withoutSpaces(name))
}

var (
// intMycology *Interest = NewInterest(string(smMycology)) // TODO: DEL?
// intCryptocurrency *Interest = NewInterest(string(smCryptocurrency))
// intCryptography *Interest = NewInterest(string(smCryptography))
)
