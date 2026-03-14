package main

import (
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
	Info        []string // Information points on a client. // TODO: ADD INFO TO ALL CLIENTS AND DISPLAY
	Projects    []*Project
	projectsSet utils.Set[string]
	//Languages []string // Calculated later
	// Parent
	company *CompanyPage
}

func (pg *Client) NameValue() string {
	return pg.Name
}

func (pg *Client) ProjectTypeInfo() projectTypeInfo {
	return professionalProjectTypeInfo{client: pg}
}

func (pg *Client) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "client", withoutSpaces(pg.Name))
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
			builder.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		}
	}
	ls := languagesFor(pg.Projects)
	if len(ls) > 0 {
		builder.WriteString("# Languages\n")
		builder.WriteString("Language | Usage Frequency[[" + withoutSpaces(pg.Name) + "#^usageFrequencies\\|*]]\n")
		// TODO: remove frequency???
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
				builder.WriteString(fmt.Sprintf("%s | %s\n", langs[name].Link(), Frequency(f).String()))
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
			providers[l].AddClient(c) // TODO: cloud services?
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
		"Created multiple "+lookup("ECS").Link()+" clusters for use in, and consuming data from, John Deere AI platforms",
		"Designed "+lookup("CUDA").Link()+" "+lookup("C").Link()+"/[C++](cv/language/Cpp) Kernels used through "+lookup("CGo").Link()+" for statistics and image processing via GPU",
		"Improved a mission-critical image manipulation API from 65% reliability to 99.9999% success rate",
		"Spearheaded implementation of an "+lookup("ECS").Link()+" cluster using an advanced "+lookup("topology").Link()+" algorithm (from a PhD thesis), achieving >100x performance gains over its original implementation on TB-scale datasets",
		"Built a "+lookup("Go").Link()+"-based compiler that transformed nested "+lookup("JSON").Link()+" instructions into machine-executable operations across clusters of "+lookup("EC2").Link()+" instances for for retrieving and manipulating geospatial agricultural data",
		"__Saved \\$18M of a \\$28M budget (64%)__ in 2024, while still increasing service stability and throughput",
		"Ensured maximum service uptime via careful design and rollout of [CI/CD](cv/subjectMatter/CI-CD) pipelines operated via self-hosted GitHub Actions runners, in conjunction with Infrastructure as Code via "+lookup("Terraform").Link(),
		"Setup monitoring, dashboards, alerting, traces, and profiling ("+lookup("Datadog").Link()+"/"+lookup("Grafana").Link()+"/"+lookup("Cloudwatch").Link()+")",
		"Responsible for educating engineers on infrastructure, codebase, domain, and best practices",
		"Utilized primarily "+lookup("Go").Link()+", "+lookup("Terraform").Link()+", "+lookup("Bash").Link()+", and "+lookup("Docker").Link()+" on "+lookup("AWS").Link()+", but also used "+lookup("Scala").Link()+", "+lookup("Github Actions").Link()+", "+lookup("C").Link()+"/[C++](cv/language/Cpp) with "+lookup("CUDA").Link()+", "+lookup("DroneCI").Link()+", "+lookup("Python").Link()+", "+lookup("Javascript").Link()+", "+lookup("Typescript").Link()+", "+lookup("Kotlin").Link()+", and more",
		"Platforms utilized:  AWS (>30 separate services), "+lookup("Datadog").Link()+", "+lookup("LogCentral").Link()+", "+lookup("Rally").Link()+", "+lookup("Azure DevOps").Link()+", "+lookup("Github").Link()+", "+lookup("DroneCI").Link()+", "+lookup("Grafana").Link()+", "+lookup("Prometheus").Link()+", "+lookup("Confluence").Link()+", and more")
	sourceAlliesClient = sourceAlliesClient.
		WithProjects(jamfProject, smallImprovementsProject, internalResumeGeneratorProject). // TODO: USE OTHERS
		WithInfo(
			"Source Allies internal projects",
			"Upgraded company internal payment gateway to a newer version of "+lookup("Java").Link()+" "+lookup("Spring").Link(),
			"Secured all company machines via "+lookup("JAMF").Link()+" to ensure the protection of company and client data",
			"Designed and created a "+lookup("Slack").Link()+" bot integration with "+lookup("Small Improvements").Link()+" to automate monthly announcements, and employee creation and completion of personal and professional goals",
			"Redesigned "+lookup("Jira").Link()+" workflows streamlining the hiring process and onboarding systems for remote coworkers",
		)
	simpsonUniversityClient = simpsonUniversityClient.
		WithProjects(simpsonUnivProject).
		WithInfo("Upgraded the University’s payment and donation gateway. Remediated resulting bugs")
	// TODO: SIMPSON COLLEGE // TODO: USE!
	critColaClient = critColaClient.
		WithProjects(CritColaProject).
		WithInfo(
			// TODO; link to discord
			"Consulted on hosting game servers and "+lookup("Discord").Link()+" bots for a large online community such that they could be deployed or destroyed, on short notice with persistent game data utilizing "+lookup("Gitlab CI").Link()+", "+lookup("Terraform").Link()+", "+lookup("Cloudflare").Link()+", and "+lookup("AWS").Link()+" ("+lookup("EC2").Link()+", "+lookup("EBS").Link()+", "+lookup("IAM").Link()+")",
			"Created a final product with a spin-up time of approximately 3 minutes",
		)
	wellAwareClient = wellAwareClient.WithProjects(WellAwareProject).
		WithInfo("Designed, created, and hosted a website for Well Aware NC, a University of North Carolina Chapel Hill affiliated nonprofit focused on the testing of well water contaminants within North Carolina")
	clarkClient = clarkClient.WithProjects(MastersDataAnalysisProject).
		WithInfo("Created programs for a student doing research for his Masters Degree in Public Health. Provided data on E.Coli samples from different waterways, the programs checked the statistical validity on different E.Coli indicating kits")
	charityClient = charityClient.WithProjects(CharityProject).
		WithInfo("Designed, created, and hosted a website for an undisclosed local charity")
	teiClient = teiClient.WithProjects(teiProjects).
		WithInfo("Internal company work for TEI") // TODO: WEBSITE?      // TODO: add projects (like NM, TX, IA, NC?)
	taeClient = taeClient.WithProjects(taePhotoImporter).
		WithInfo("Internal company work for Talley Associates of Engineering")
	mafcClient = mafcClient.WithProjects(mafcProjects).
		WithInfo("Lifeguarding work, both at the indoor and outdoor pools of the fitness center") // TODO: Indoor and outdoor pool?
	// TODO: links in text for everything below this
	arrowNailClient = arrowNailClient.WithProjects(ArrowNailProject).
		WithInfo("Used "+lookup("Serverless").Link()+" services on "+lookup("AWS").Link()+" to support a "+lookup("React").Link()+" geospatial web app, Node API via "+lookup("lambda").Link()+" functions, and an "+lookup("Aurora").Link()+" database. Provided IT and Systems Administration Support",
			"Set up [CI/CD](cv/subjectMatter/CI-CD) pipeline in "+lookup("Gitlab CI").Link()+" to make future deployments seamless")
	wildlifeRClient = wildlifeRClient.WithProjects(WildlifeRProject).
		WithInfo("Utilized spatiotemporal data for wildlife in a specified area over a specified date range in order to produce population density maps")
}
