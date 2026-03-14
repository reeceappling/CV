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
	company        *CompanyPage
	miscSkills     utils.Set[string]        // TODO: MISC SKILLS??????
	SubjectMatters utils.Set[SubjectMatter] // TODO: DISPLAY THIS??? // TODO: maybe dont have this on here
}

func (pg *Position) NameValue() string { // The name to be used on the link. NOT the URL
	return pg.Name // TODO: or mapName?
}

func (pg *Position) EntryType() string {
	return "Position"
}

func (pg *Position) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "position", withoutSpaces(pg.MapName()))
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
func (pg *Position) WithSubjectMatters(sms ...SubjectMatter) *Position {
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
		End:            nil,
		miscSkills:     map[string]struct{}{},
		SubjectMatters: map[SubjectMatter]struct{}{},
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
	builder.WriteString(frontmatterFor(pg.Name))
	// TODO: Start/end
	// Company and clients
	builder.WriteString(fmt.Sprintf("Company: %s\n", pg.company.Link()))
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
			builder.WriteString(fmt.Sprintf("- %s\n", miscSkills[skill].Link()))
		}
	}
	// SubjectMatters
	if len(pg.SubjectMatters) > 0 {
		builder.WriteString("# Related Subject Matters\n")
		for sm, _ := range pg.SubjectMatters {
			builder.WriteString(fmt.Sprintf("- %s\n", sm.Link()))
		}
	}

	// ALL LOWEST
	builder.WriteString(bytesForAll(pg, false))
	return []byte(builder.String())
}

// PROJECT GOES ON CLIENT AND POSITION. SHOULD GO ON CLIENT FIRST
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
	// TODO: MISC SKILLS?
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

// TODO; ensure all professional projects are added to positions
func initPositionsAfterProjects() {
	// TODO: ANY MISC SKILLS
	positionLifeguard = NewPosition("Lifeguard", 1, 2013, utils.Pointer(6), utils.Pointer(2014)). // TODO: ensure dates are right
		WithSubjectMatters(smFirstAid)
	positionSeniorLifeguard = NewPosition("Senior Lifeguard", 6, 2014, utils.Pointer(8), utils.Pointer(2016)). // TODO: ensure dates are right
		WithSubjectMatters(smFirstAid)
	positionTAE = NewPosition("Civil Structural Engineer and Tower Climber", 5, 2017, utils.Pointer(3), utils.Pointer(2018)). // TODO: ensure dates are right
		WithMiscSkills("Excel", "Climbing", "AutoDesk Inventor", "AutoCAD", "Autodesk Revit", "Drafting").
		WithSubjectMatters(smCivilEngineering, smStructuralEngineering, smStatics)
	positionTEI = NewPosition("Cell Tower Inspector and Tower Climber", 1, 2019, utils.Pointer(3), utils.Pointer(2020)). // TODO: ensure dates are right
		WithMiscSkills("Climbing", "Drafting").
		WithSubjectMatters(smCivilEngineering, smStructuralEngineering, smStatics)
	freelancePosition = NewPosition("Software Engineer", 1, 2012, utils.Pointer(5), utils.Pointer(2022)).
		WithSubjectMatters(smFullStack).
		WithProjects(WellAwareProject, CharityProject, CritColaProject, MastersDataAnalysisProject, WildlifeRProject, ArrowNailProject)
	sai1 = NewPosition("Software Engineer", 5, 2022, utils.Pointer(6), utils.Pointer(2023)).
		WithSubjectMatters(smFullStack).
		WithProjects(rasterRenderProject, explorerProject, supportProject, scudsProject, simpsonUnivProject, jamfProject, smallImprovementsProject, internalResumeGeneratorProject) // TODO: MOVE PROJECTS AROUND
	sai2 = NewPosition("Senior Software Engineer", 6, 2023, utils.Pointer(2), utils.Pointer(2025)).
		WithSubjectMatters(smBackend).
		WithProjects(tileGenProject, polygonBuilderProject, GhaRunnersProject, statsProject, ufoProject, renderProject)
	sai3 = NewPosition("Senior Software Engineer and Tech Lead", 2, 2025, nil, nil).
		WithSubjectMatters(smBackend).
		WithProjects(ogreProject, billingProject, goweProject, wqdbProject)

}
