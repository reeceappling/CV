package main

import (
	"appli.ng/cv/generator/tags"
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var clients = map[string]*Client{}
var clientsOrder = []string{}

type Client struct {
	Name        string
	Description string   // TODO: ADD THIS?
	Info        []string // Information points on a client. // TODO: ADD INFO TO ALL CLIENTS AND DISPLAY
	Projects    []*Project
	projectsSet utils.Set[string]
	//Languages []string // Calculated later
	// Parent
	company *CompanyPage
	Tags    tags.Field
}

func (pg *Client) Title() string {
	if pg == nil {
		return "NO_NAME"
	}
	return pg.Name
}
func (pg *Client) WithTags(tags ...tags.Tag) *Client {
	pg.Tags = append(pg.Tags, tags...)
	return pg
}

func (pg *Client) NameValue() string {
	return pg.Name
}

func (pg *Client) ProjectTypeInfo() projectTypeInfo {
	return professionalProjectTypeInfo{client: pg}
}

const noLinkText = "NO_LINK"

func (pg *Client) Dst() string {
	if pg == nil {
		return noLinkText
	}
	return dstFor("cv", "client", withoutSpaces(pg.Name))
}
func (pg *Client) EntryType() string {
	return "Client"
}

func NewClient(name string) *Client {
	out := &Client{Name: name, Projects: []*Project{}, Info: nil, projectsSet: utils.Set[string]{}}
	if _, exists := clients[name]; exists {
		panic("client already exists")
	}
	clients[name] = out
	addLinkable(strings.ToLower(name), out)
	clientsOrder = append(clientsOrder, name)
	return out
}

func (pg *Client) WithInfo(info ...string) *Client {
	pg.Info = info
	return pg
}

func (pg *Client) Bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(frontmatterFor(pg.Name, "Client"))
	if pg.Info != nil && len(pg.Info) > 0 {
		builder.WriteString("# Notable Information\n")
		for _, info := range pg.Info {
			builder.WriteString(fmt.Sprintf("- %s\n", info))
		}
	}
	if pg.Projects != nil && len(pg.Projects) > 0 {
		builder.WriteString("# Projects\n")
		for _, proj := range pg.Projects {
			builder.WriteString(fmt.Sprintf("- %s\n", Link(proj)))
		}
	}
	ls := languagesFor(pg.Projects)
	if len(ls) > 0 {
		builder.WriteString("# Languages\n")
		builder.WriteString("Language | Usage Frequency[[" + withoutSpaces(pg.Name) + "#^usageFrequencies\\|*]]\n")
		builder.WriteString(":-- | :--\n")
		freqs := make(map[Frequency][]string, 6)
		for name, freq := range ls {
			if freqs[freq] == nil {
				freqs[freq] = []string{name}
			} else {
				freqs[freq] = append(freqs[freq], name)
			}
		}
		for f := 5; f >= 0; f-- {
			for _, name := range freqs[Frequency(f)] {
				builder.WriteString(fmt.Sprintf("%s | %s\n", Link(langs[name]), Frequency(f).String()))
			}
		}
	}
	// ALL LOWEST
	builder.WriteString(bytesForAll(pg, true))
	return []byte(builder.String())
}

func projectsFor(cs []*Client) []*Project {
	temp := map[string]*Project{}
	for _, cli := range cs {
		for _, pr := range cli.Projects {
			temp[pr.Name] = pr
		}
	}
	return slices.Collect(maps.Values(temp))
}

func (c *Client) WithProjects(projs ...*Project) *Client {
	for _, proj := range projs {
		if !c.projectsSet.Contains(proj.Name) {
			c.projectsSet.Add(proj.Name)
			c.Projects = append(c.Projects, proj)
		}
		proj.TypeInfo = proj.TypeInfo.setClient(c)
		for lang, _ := range proj.Languages {
			langs[lang].AddClient(c)
		}
		for l := range maps.Keys(proj.Dbs) {
			dbs[l].AddClient(c)
		}
		for l := range maps.Keys(proj.Caches) {
			caches[l].AddClient(c)
		}
		for l := range maps.Keys(proj.Technologies) {
			techs[l].AddClient(c)
		}
		for l := range maps.Keys(proj.CloudProviders) {
			providers[l].AddClient(c)
		}
		for l := range maps.Keys(proj.Platforms) {
			platforms[l].AddClient(c)
		}
	}

	return c
}

func (pg *Client) GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
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
	}
	return
}

func initClientsLast() {

}

var (
	jdClient                = NewClient("John Deere")
	simpsonUniversityClient = NewClient("Simpson University")
	sourceAlliesClient      = NewClient("Source Allies")
	critColaClient          = NewClient("CritCola")
	arrowNailClient         = NewClient("ArrowNail LLC")
	wellAwareClient         = NewClient("Well Aware NC")
	clarkClient             = NewClient("Chapel Hill Masters Student in Public Health")
	wildlifeRClient         = NewClient("Wildlife mapping with R")
	charityClient           = NewClient("Undisclosed Charity")
	teiClient               = NewClient("TEI")
	taeClient               = NewClient("Talley Associates of Engineering")
	mafcClient              = NewClient("Monroe Aquatics and Fitness Center")
)

func initClientsAfterProjectsComplete() {
	jdClient = jdClient.
		WithProjects(polygonBuilderProject, ogreProject, renderProject, statsProject, billingProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, tileGenProject, GhaRunnersProject).WithInfo(
		"Most senior consulting engineer on a high-performance, global-scale, team of 4-8 at John Deere’s Intelligent Solutions Group; responsible for architecture, implementation, testing, optimization, and support of a complex set of diverse cloud services utilizing geospatiotemporal agribusiness data",
		"Architected, implemented, and maintained cloud infrastructure and services for ingest, distributed processing, storage, manipulation, and retrieval of data for, datastores totalling over 50PB",
		"Created multiple "+Link(lookup("ECS"))+" clusters for use in, and consuming data from, John Deere "+Link(lookup("AI"))+" platforms",
		"Designed "+Link(lookup("CUDA"))+" "+Link(lookup("C"))+"/[C++](cv/language/Cpp) Kernels used through "+Link(lookup("CGo"))+" for statistics and image processing via GPU",
		"Improved a mission-critical image manipulation API from 65% reliability to 99.9999% success rate",
		"Spearheaded implementation of an "+Link(lookup("ECS"))+" cluster using an advanced "+Link(lookup("topology"))+" algorithm (from a PhD thesis), achieving >100x performance gains over its original implementation on TB-scale datasets",
		"Built a "+Link(lookup("Go"))+"-based compiler that transformed nested "+Link(lookup("JSON"))+" instructions into machine-executable operations across clusters of "+Link(lookup("EC2"))+" instances for for retrieving and manipulating geospatial agricultural data",
		"__Saved \\$18M of a \\$28M budget (64%)__ in 2024, while still increasing service stability and throughput",
		"Ensured maximum service uptime via careful design and rollout of "+Link(lookup("CI-CD"))+" pipelines operated via self-hosted GitHub Actions runners, in conjunction with Infrastructure as Code via "+Link(lookup("Terraform")),
		"Setup monitoring, dashboards, alerting, traces, and profiling ("+Link(lookup("Datadog"))+"/"+Link(lookup("Grafana"))+"/"+Link(lookup("Cloudwatch"))+")",
		"Responsible for educating engineers on infrastructure, codebase, domain, and best practices",
		"Utilized primarily "+Link(lookup("Go"))+", "+Link(lookup("Terraform"))+", "+Link(lookup("Bash"))+", and "+Link(lookup("Docker"))+" on "+Link(lookup("AWS"))+", but also used "+Link(lookup("Scala"))+", "+Link(lookup("Github Actions"))+", "+Link(lookup("C"))+"/[C++](cv/language/Cpp) with "+Link(lookup("CUDA"))+", "+Link(lookup("DroneCI"))+", "+Link(lookup("Python"))+", "+Link(lookup("Javascript"))+", "+Link(lookup("Typescript"))+", "+Link(lookup("Kotlin"))+", and more",
		"Platforms utilized:  AWS (>30 separate services), "+Link(lookup("Datadog"))+", "+Link(lookup("LogCentral"))+", "+Link(lookup("Rally"))+", "+Link(lookup("Azure DevOps"))+", "+Link(lookup("Github"))+", "+Link(lookup("DroneCI"))+", "+Link(lookup("Grafana"))+", "+Link(lookup("Prometheus"))+", "+Link(lookup("Confluence"))+", and more").
		WithTags(tags.AgTech)
	sourceAlliesClient = sourceAlliesClient.
		WithProjects(jamfProject, smallImprovementsProject, internalResumeGeneratorProject). // TODO: USE OTHERS
		WithInfo(
			"Source Allies internal projects",
			"Upgraded company internal payment gateway to a newer version of "+Link(lookup("Java"))+" "+Link(lookup("Spring")),
			"Secured all company machines via "+Link(lookup("JAMF"))+" to ensure the protection of company and client data",
			"Designed and created a "+Link(lookup("Slack"))+" bot integration with "+Link(lookup("Small Improvements"))+" to automate monthly announcements, and employee creation and completion of personal and professional goals",
			"Redesigned "+Link(lookup("Jira"))+" workflows streamlining the hiring process and onboarding systems for remote coworkers",
		).
		WithTags(tags.Consulting)
	simpsonUniversityClient = simpsonUniversityClient.
		WithProjects(simpsonUnivProject).
		WithInfo("Upgraded the University’s payment and donation gateway. Remediated resulting bugs").
		WithTags(tags.Education)
	critColaClient = critColaClient.
		WithProjects(CritColaProject).
		WithInfo(
			"Consulted on hosting game servers and "+Link(lookup("Discord"))+" bots for a large online community such that they could be deployed or destroyed, on short notice with persistent game data utilizing "+Link(lookup("Gitlab CI"))+", "+Link(lookup("Terraform"))+", "+Link(lookup("Cloudflare"))+", and "+Link(lookup("AWS"))+" ("+Link(lookup("EC2"))+", "+Link(lookup("EBS"))+", "+Link(lookup("IAM"))+")",
			"Created a final product with a spin-up time of approximately 3 minutes",
		).WithTags(tags.Entertainment)
	wellAwareClient = wellAwareClient.WithProjects(WellAwareProject).
		WithInfo("Designed, created, and hosted a website for Well Aware NC, a University of North Carolina Chapel Hill affiliated nonprofit focused on the testing of well water contaminants within North Carolina").
		WithTags(tags.PublicHealth, tags.Education)
	clarkClient = clarkClient.WithProjects(MastersDataAnalysisProject).
		WithInfo("Created programs for a student doing research for his Masters Degree in Public Health. Provided data on E.Coli samples from different waterways, the programs checked the statistical validity on different E.Coli indicating kits").
		WithTags(tags.PublicHealth, tags.Education)
	charityClient = charityClient.WithProjects(CharityProject).
		WithInfo("Designed, created, and hosted a website for an undisclosed local charity").
		WithTags(tags.Charity)
	teiClient = teiClient.WithProjects(teiProjects).
		WithInfo("Internal company work for TEI"). // TODO: WEBSITE?      // TODO: add projects (like NM, TX, IA, NC?)
		WithTags(tags.Telecom)
	taeClient = taeClient.WithProjects(taePhotoImporter).
		WithInfo("Internal company work for Talley Associates of Engineering").
		WithTags(tags.Telecom)
	mafcClient = mafcClient.WithProjects(mafcProjects).
		WithInfo("Lifeguarding work, both at the indoor and outdoor pools of the fitness center"). // TODO: Indoor and outdoor pool?
		WithTags(tags.FirstAid)
	arrowNailClient = arrowNailClient.WithProjects(ArrowNailProject).
		WithInfo("Used "+Link(lookup("Serverless"))+" services on "+Link(lookup("AWS"))+" to support a "+Link(lookup("React"))+" geospatial web app, Node API via "+Link(lookup("lambda"))+" functions, and an "+Link(lookup("Aurora"))+" database. Provided IT and Systems Administration Support",
			"Set up "+Link(lookup("CI-CD"))+" pipeline in "+Link(lookup("Gitlab CI"))+" to make future deployments seamless").
		WithTags(tags.Construction)
	wildlifeRClient = wildlifeRClient.WithProjects(WildlifeRProject).
		WithInfo("Utilized spatiotemporal data for wildlife in a specified area over a specified date range in order to produce population density maps").
		WithTags(tags.Ecology)
}
