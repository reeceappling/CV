package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"strings"
)

var schools = map[string]*SchoolPage{}

type SchoolPage struct {
	Name             string
	Info             *string // use this when the company can't be outright named!
	Degrees          []*Degree
	Projects         []*Project
	Extracurriculars []Extracurricular
	SubjectMatters   utils.Set[SubjectMatter] // TODO: DISPLAY THIS???
}

func NewSchool(name string) *SchoolPage {
	out := &SchoolPage{
		Name:             name,
		Info:             nil,
		Degrees:          []*Degree{},
		Projects:         []*Project{},
		Extracurriculars: []Extracurricular{},
		SubjectMatters:   map[SubjectMatter]struct{}{},
	}
	schools[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}
func (pg *SchoolPage) WithSummary(info string) *SchoolPage {
	pg.Info = &info
	return pg
}
func (pg *SchoolPage) WithSubjectMatters(sms ...SubjectMatter) *SchoolPage {
	if pg == nil {
		return pg
	}
	for _, sm := range sms {
		if sm == smFullStack {
			pg.SubjectMatters.Add(smFrontend, smBackend)
		}
		pg.SubjectMatters.Add(sm)
	}

	return pg
}
func (pg *SchoolPage) ProjectTypeInfo() projectTypeInfo {
	return schoolProjectTypeInfo{school: pg}
}

func (pg *SchoolPage) Dst() string {
	return dstFor("cv", "school", withoutSpaces(pg.Name))
}

func (pg *SchoolPage) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
}

func (pg *SchoolPage) EntryType() string {
	return "School"
}

func (sp *SchoolPage) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString(frontmatterFor(sp.Name, "School"))
	// Write all degrees
	b.WriteString("# Degrees\n")
	for i, deg := range sp.Degrees {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(deg.String())
	}
	b.WriteString("\n")
	// Write all projects
	b.WriteString("# Projects\n")
	for _, proj := range sp.Projects {
		b.WriteString(fmt.Sprintf("- %s\n", Link(proj)))
	}
	// Write all extracurriculars
	b.WriteString("# Extracurriculars and positions held\n")
	for _, ec := range sp.Extracurriculars {
		b.WriteString(fmt.Sprintf("%s - %s\n", ec.Name, ec.Summary))
		if len(ec.Positions) != 0 {
			for _, pos := range ec.Positions {
				b.WriteString(fmt.Sprintf("- %s", pos.Title))
				if ec.Summary != "" {
					b.WriteString(fmt.Sprintf(" - %s\n", pos.notes))
				}
			}
		}
	}
	// SubjectMatters
	if len(sp.SubjectMatters) > 0 {
		b.WriteString("# Related Subject Matters\n")
		for sm, _ := range sp.SubjectMatters {
			b.WriteString(fmt.Sprintf("- %s\n", Link(sm)))
		}
	}

	return []byte(b.String())
}

func (sp *SchoolPage) withProjects(projs ...*Project) *SchoolPage {
	sp.Projects = append(sp.Projects, projs...)
	for _, proj := range projs {
		proj.TypeInfo = proj.TypeInfo.setSchool(sp)
		for lang, _ := range proj.Languages {
			langs[lang].AddSchool(sp)
		}
		for l := range maps.Keys(proj.Dbs) {
			dbs[l].AddSchool(sp)
		}
		for l := range maps.Keys(proj.Caches) {
			caches[l].AddSchool(sp)
		}
		for l := range maps.Keys(proj.Technologies) {
			techs[l].AddSchool(sp)
		}
		for l := range maps.Keys(proj.CloudProviders) {
			providers[l].AddSchool(sp) // TODO: cloud services?
		}
		for l := range maps.Keys(proj.Platforms) {
			platforms[l].AddSchool(sp)
		}
	}

	return sp
}

func (sp *SchoolPage) WithDegree(deg *Degree) *SchoolPage {
	sp.Degrees = append(sp.Degrees, deg)
	return sp
}

func (sp *SchoolPage) WithExtracurriculars(ex ...Extracurricular) *SchoolPage {
	sp.Extracurriculars = append(sp.Extracurriculars, ex...)
	return sp
}

type Extracurricular struct {
	Name      string
	Summary   string
	Positions []EcPosition
}

func NewExtracurricular(name string, summary string, positions ...EcPosition) Extracurricular {
	return Extracurricular{
		Name:      name,
		Summary:   summary,
		Positions: positions,
	}
}

type EcPosition struct {
	Title string
	notes string
}

func NewEcPosition(title, notes string) EcPosition {
	return EcPosition{title, notes}
}

var (
	schoolCata = NewSchool("Central Academy of Technology and Arts").
			WithSummary("Magnet High School, Engineering").
			WithSubjectMatters(smEducation, smStatics, smElectronics).
			WithDegree(DegHS).
			WithExtracurriculars(
			NewExtracurricular("Soccer", fixmeLink,
				NewEcPosition("Varsity", "Sophomore-Senior year"),
				NewEcPosition("Junior Varsity", "Freshman year"),
			),
			NewExtracurricular("Swimming", fixmeLink,
				NewEcPosition("Varsity", "Senior year"),
			),
			NewExtracurricular("Robotics club", "FRC team 3720",
				NewEcPosition("Team Captain", "Junior-Senior year"),
				NewEcPosition("Driver", "Junior-Senior year"),
				NewEcPosition("Member", "Freshman-Senior year"),
				NewEcPosition("Lead programmer", "Junior-Senior year"),
				NewEcPosition("Programmer", "Freshman year"),
			),
			NewExtracurricular("Student Council", fixmeLink,
				NewEcPosition("Treasurer", "Senior year"),
			),
			NewExtracurricular("National Honor Society (NHS)", fixmeLink),
			NewExtracurricular("Beta club", fixmeLink),
		)
	schoolNCSU = NewSchool("North Carolina State University").
			WithSummary("Undergraduate studies").
			WithDegree(DegNE).
			WithDegree(DegMath).
			WithSubjectMatters(smNuclearEngineering, smParticlePhysics, smFluidMechanics, smThermodynamics, smEducation, smStatics, smElectronics).
			WithExtracurriculars(
			NewExtracurricular("American Nuclear Society", fixmeLink),
			NewExtracurricular("88.1 WKNC FM HD1 Raleigh", fixmeLink,
				NewEcPosition("DJ", "Sophomore-Junior years"),
			),
			NewExtracurricular("Wolftrax music group", fixmeLink),
		)
)

type Degree struct {
	Type      string
	TypeShort string
	Area      string
}

func (deg Degree) String() string {
	return deg.Type + " in " + deg.Area
}
func (deg Degree) StringShort() string {
	return deg.TypeShort + " " + deg.Area
}

var DegNE = NewDegree("Bachelor of Science", "B.S.", "Nuclear Engineering")
var DegMath = NewDegree("Minor", "Minor", "Mathematics")
var DegHS = NewDegree("High School Diploma", "HS Diploma", "Engineering Academy")

func NewDegree(Type, TypeShort, Area string) *Degree {
	return &Degree{
		Type:      Type,
		TypeShort: TypeShort,
		Area:      Area,
	}
}
