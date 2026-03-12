package main

var miscSkills = map[string]*MiscSkill{}

type MiscSkill struct {
	Name string
	*tracked
}

func (pg *MiscSkill) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "miscSkills", withoutSpaces(pg.Name)) // TODO: plural??
}

func NewSkill(name string) *MiscSkill {
	out := &MiscSkill{
		Name:    name,
		tracked: newTracked(),
	}
	miscSkills[name] = out
	return out
}
