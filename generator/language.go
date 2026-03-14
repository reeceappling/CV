package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"strings"
)

var langs = map[string]*LanguagePage{}

type LanguagePage struct {
	Name      string
	Companies map[string]*CompanyPage
	Clients   map[string]*Client
	Schools   map[string]*SchoolPage
	Projects  map[string]Frequency
	Positions map[string]*Position
	SubjectMattersField
}

func (pg *LanguagePage) EntryType() string {
	return "Language"
}

func (pg *LanguagePage) WithSubjectMatters(sms ...SubjectMatter) *LanguagePage {
	if pg == nil {
		return nil
	}
	pg.SubjectMatters = withSubjectMatters(pg, pg.SubjectMatters, sms...)
	return pg
}

func (pg *LanguagePage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "language", withoutSpaces(pg.Name))
}

func (pg *LanguagePage) Bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(frontmatterFor(pg.Name, "Language"))
	if pg.Companies != nil && len(pg.Companies) > 0 {
		builder.WriteString("# Companies\n")
		for companyName := range pg.Companies {

			builder.WriteString(fmt.Sprintf("- %s\n", companies[companyName].Link()))
		}
	}
	if pg.Clients != nil && len(pg.Clients) > 0 {
		builder.WriteString("# Clients\n")
		for name := range pg.Clients {
			builder.WriteString(fmt.Sprintf("- %s\n", clients[name].Link()))
		}
	}
	if pg.Projects != nil && len(pg.Projects) > 0 {
		builder.WriteString("# Projects\n")
		builder.WriteString("Project | Usage Frequency[[" + withoutSpaces(pg.Name) + "#^usageFrequencies\\|*]] | Project Type\n")
		builder.WriteString(":-- | :--: | :--\n")
		freqs := make(map[Frequency][]string, 6)
		for name, freq := range pg.Projects {
			if freqs[freq] == nil {
				freqs[freq] = []string{name}
			} else {
				freqs[freq] = append(freqs[freq], name)
			}
		}
		for f := 5; f >= 0; f-- {
			for _, name := range freqs[Frequency(f)] {
				projTypeStr := "unknown"
				v := projects[name].TypeInfo.Type()
				switch v {
				case projectTypePersonal:
					projTypeStr = "Personal"
				case projectTypeProfessional:
					projTypeStr = "Professional"
				case projectTypeSchool:
					projTypeStr = "Coursework-related"
				default:
					panic("unknown project type: " + v)
				}

				builder.WriteString(fmt.Sprintf("%s | %s | %s\n", projects[name].Link(), Frequency(f).String(), projTypeStr))
			}
		}
	}
	builder.WriteString(definitionsArea())
	builder.WriteString(usageFrequencyDefinitions())
	return []byte(builder.String())
}
func definitionsArea() string {
	return "### Definitions\n"
}
func usageFrequencyDefinitions() string {
	return "#### * Usage Frequencies\n" +
		"```text\n" +
		"Minimal: Very little use during a single project\n" +
		"Rarely: Sparing use during a single project\n" +
		"Some: Occasional use during a single project\n" +
		"Regularly: Not the dominant language used during a single project, but still commonly used\n" +
		"Often:  Not the most often used language in a single project, but still very frequently utilized\n" +
		"Extensively: The dominant, most often used language in a single project\n" +
		"```\n^usageFrequencies"
}

func NewLanguage(name string) *LanguagePage {
	if _, exists := langs[name]; exists {
		panic("tried to create language twice: " + name)
	}
	out := &LanguagePage{
		Name:                name,
		Companies:           map[string]*CompanyPage{},
		Clients:             map[string]*Client{},
		Projects:            map[string]Frequency{},
		Schools:             map[string]*SchoolPage{},
		Positions:           map[string]*Position{},
		SubjectMattersField: SubjectMattersField{utils.Set[SubjectMatter]{}},
	}
	langs[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}

func (lp *LanguagePage) AddClient(client *Client) *LanguagePage {
	lp.Clients[client.Name] = client
	return lp
}
func (lp *LanguagePage) AddSchool(school *SchoolPage) *LanguagePage {
	lp.Schools[school.Name] = school
	return lp
}
func (lp *LanguagePage) AddProject(proj *Project) *LanguagePage {
	lp.Projects[proj.Name] = proj.Languages[lp.Name]
	return lp
}
func (lp *LanguagePage) AddCompany(comp *CompanyPage) *LanguagePage {
	lp.Companies[comp.Name] = comp
	return lp
}
func (lp *LanguagePage) AddPosition(pos *Position) *LanguagePage {
	lp.Positions[pos.Name] = pos
	return lp
}

func setupLanguageSubjectMatters() {
	NewLanguage("Docker").
		WithSubjectMatters(smContainerization)
	NewLanguage("Docker Compose").
		WithSubjectMatters(smContainerization)
	NewLanguage("Terraform").
		WithSubjectMatters(smIAC)
	NewLanguage("Html").
		WithSubjectMatters(smFrontend)
	NewLanguage("CSS").
		WithSubjectMatters(smFrontend)
	NewLanguage("Javascript").
		WithSubjectMatters(smFullStack)
}
