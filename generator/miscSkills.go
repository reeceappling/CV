package main

var miscSkills = map[string]*MiscSkill{}

type MiscSkill struct {
	Name string
	*tracked
}

func (pg *MiscSkill) Dst() string {
	return dstFor("cv", "miscSkill", withoutSpaces(pg.Name))
}

func (pg *MiscSkill) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
}

func (pg *MiscSkill) EntryType() string {
	return "Misc Skill"
}

func NewSkill(name string) *MiscSkill {
	out := &MiscSkill{
		Name:    name,
		tracked: newTracked(),
	}
	miscSkills[name] = out
	return out
}
