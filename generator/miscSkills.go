package main

import "fmt"

var miscSkills = map[string]*MiscSkill{}

type MiscSkill struct {
	Name string
	*tracked
}

func (pg *MiscSkill) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](miscSkills/%s)", pg.Name, withoutSpaces(pg.Name))
}

func NewSkill(name string) *MiscSkill {
	out := &MiscSkill{
		Name:    name,
		tracked: newTracked(),
	}
	miscSkills[name] = out
	return out
}
