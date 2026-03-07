package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"strings"
)

var positionsMap = map[string]*Position{} // Map of companyName+positionName to position

type Position struct {
	Name  string
	Start monthYr
	End   *monthYr // None == current
	//Clients  []*Client // clients may contain projects which are outside of this position
	Projects map[string]*Project
	// Parent
	company    *CompanyPage
	miscSkills utils.Set[string]
	// TODO: MISC SKILLS??????
}

func (pg *Position) NameValue() string {
	return pg.Name
}

func (pg *Position) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](position/%s)", pg.Name, withoutSpaces(pg.MapName()))
}
func (pg *Position) WithMiscSkills(skills ...string) *Position {
	if pg == nil {
		return pg
	}
	pg.miscSkills.Add(skills...)
	for _, skill := range skills {
		_, exists := miscSkills[skill]
		if !exists {
			_ = NewSkill(skill)
		}
	}
	return pg
}

func NewPosition(name string, startMo, startYr int, endMo, endYr *int) *Position {
	if _, exists := companies[name]; exists {
		panic("company already exists")
	}
	out := &Position{
		Name: name,
		Start: monthYr{
			Month: startMo,
			Year:  startYr,
		},
		End:        nil,
		miscSkills: map[string]struct{}{},
	}
	if endMo != nil && endYr != nil {
		out.End = &monthYr{
			Month: *endMo,
			Year:  *endYr,
		}
	}
	// NOTE: not added to the map here. Added to the map when added to a company!
	return out
}

func (pg *Position) Bytes() []byte {
	builder := strings.Builder{}
	// TODO: NAME
	// TODO: Start/end
	// Company and clients
	builder.WriteString(fmt.Sprintf("Company: [%s](company/%s)\n", pg.company.Name, withoutSpaces(pg.company.Name)))
	if len(pg.Projects) > 0 {
		builder.WriteString("# Clients and Projects\n")
		builder.WriteString("Client | Project\n")
		builder.WriteString(":-- | :--\n")
		for _, pr := range pg.Projects {
			builder.WriteString(fmt.Sprintf("%s | %s\n", pr.TypeInfo.getClient().Link(), pr.Link()))
		}
	} else {
		builder.WriteString("NO CLIENTS FOUND. FIXME\n ")
	}
	// Misc skills
	if len(pg.miscSkills) > 0 {
		builder.WriteString("# Misc Skills\n")
		for skill, _ := range pg.miscSkills {
			builder.WriteString(fmt.Sprintf("- [%s](miscSkill/%s)\n", skill, withoutSpaces(skill))) // TODO: FIX ME!
		}
	}

	//// Clients
	//if len(pg.Projects) > 0 { // TODO: CLIENTS ARE CURRENTLY BROKEN
	//
	//	tempC := make([][]string, len(pg.Clients))
	//	for i, cl := range pg.Clients {
	//		tempC[i] = make([]string, len(cl.Projects)+1)
	//		tempC[i][1] = fmt.Sprintf("- %s\n", cl.Link())
	//		for _, pr := range cl.Projects {
	//			builder.WriteString(fmt.Sprintf("%s | %s\n", cl.Link(), pr.Link()))
	//		}
	//	}
	//
	//}
	// ALL LOWEST
	builder.WriteString(bytesForAll(pg, false))
	return []byte(builder.String())
}

// TODO: PROJECT GOES ON CLIENT AND POSITION!!!!! SHOULD GO ON CLIENT FIRST!!!!!
func (pg *Position) MapName() string {
	if pg.company == nil {
		panic("no company name!")
	}
	return pg.company.Name + " " + pg.Name
}

func (pg *Position) WithProjects(projs ...*Project) *Position {
	pg.Projects = map[string]*Project{}
	for _, proj := range projs {
		pg.Projects[proj.Name] = proj
		itemCaches, itemDbs, itemLanguages, itemPlatforms, itemProviders, itemServices, itemTechs := proj.GetAllLowest()
		for _, item := range itemCaches {
			item.AddPosition(pg)
		}
		for _, item := range itemDbs {
			item.AddPosition(pg)
		}
		for langName, _ := range itemLanguages {
			langs[langName].AddPosition(pg)
		}
		for _, item := range itemPlatforms {
			item.AddPosition(pg)
		}
		for _, item := range itemProviders {
			item.AddPosition(pg)
		}
		for _, provSvcs := range itemServices {
			for _, svc := range provSvcs {
				svc.AddPosition(pg)
			}
		}
		for _, item := range itemTechs {
			item.AddPosition(pg)
		}
		pg.miscSkills.Add(proj.MiscSkills.AsSlice()...)
	}
	// TODO: put position on clients too?
	return pg
}

//func (pg *Position) Clients() map[string]*Client {
//	out := make(map[string]*Client)
//	for _, proj := range pg.Projects {
//		out[proj.client.Name] = proj.client
//	}
//	return out
//}

func initLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
	outCaches = map[string]*CachePage{}
	outDbs = map[string]*DbPage{}
	outLanguages = map[string]Frequency{}
	outPlatforms = map[string]*PlatformPage{}
	outProviders = map[string]*CloudProviderPage{}
	outServices = map[string]map[string]*CloudServicePage{}
	outTechnologies = map[string]*TechologyPage{}
	// TODO: MISC SKILLS
	return
}

func (pg *Position) GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
	outCaches, outDbs, outLanguages, outPlatforms, outProviders, outServices, outTechnologies = initLowest()
	for _, proj := range pg.Projects {
		tempC, tempD, tempL, tempP, tempPr, tempS, tempT := proj.GetAllLowest()
		for key, val := range tempC {
			outCaches[key] = val
		}
		for key, val := range tempD {
			outDbs[key] = val
		}
		for key, val := range tempL {
			if freq, exists := outLanguages[key]; exists {
				// only replace if freq is better!
				if val > freq {
					outLanguages[key] = val
				}
			}
		}
		for key, val := range tempP {
			outPlatforms[key] = val
		}
		for key, val := range tempPr {
			outProviders[key] = val
		}
		for key, val := range tempS {
			outServices[key] = val
		}
		for key, val := range tempT {
			outTechnologies[key] = val
		}
		// TODO: MISC SKILLS
	}
	return
}

var (
	sai1, sai2, sai3, freelancePosition, positionTEI, positionTAE, positionSeniorLifeguard, positionLifeguard *Position
)

func initPositionsAfterProjects() {
	// TODO: ANY MISC SKILLS
	// TODO: BEFORE OR AFTER CLIENTS?
	positionLifeguard = NewPosition("Lifeguard", 1, 2013, utils.Pointer(6), utils.Pointer(2014))                              // TODO: ensure dates are right
	positionSeniorLifeguard = NewPosition("Senior Lifeguard", 6, 2014, utils.Pointer(8), utils.Pointer(2016))                 // TODO: ensure dates are right
	positionTAE = NewPosition("Civil Structural Engineer and Tower Climber", 5, 2017, utils.Pointer(3), utils.Pointer(2018)). // TODO: ensure dates are right
																	WithMiscSkills("Excel", "Climbing Cell Towers")
	positionTEI = NewPosition("Tower Climber", 1, 2019, utils.Pointer(3), utils.Pointer(2020)). // TODO: ensure dates are right
													WithMiscSkills("Climbing Cell Towers")
	freelancePosition = NewPosition("Software Engineer", 1, 2012, utils.Pointer(5), utils.Pointer(2022)).
		WithProjects(WellAwareProject, CharityProject, CritColaProject, MastersDataAnalysisProject)
	sai1 = NewPosition("Software Engineer", 5, 2022, utils.Pointer(6), utils.Pointer(2023)).
		WithProjects(tileGenProject, renderProject, statsProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, saiCollegeProject) // TODO: MOVE PROJECTS AROUND
	sai2 = NewPosition("Senior Software Engineer", 6, 2023, utils.Pointer(2), utils.Pointer(2025)).
		WithProjects(polygonBuilderProject, GhaRunnersProject) // TODO: FIXME! polygonBuilder
	sai3 = NewPosition("Senior Software Engineer and Tech Lead", 2, 2025, nil, nil).
		WithProjects(ogreProject, billingProject) // TODO: FIXME! OGRE/Billing builder

}

var ()
