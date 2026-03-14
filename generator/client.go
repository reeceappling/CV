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

func NewClient(name string, info ...string) *Client {
	out := &Client{Name: name, Projects: []*Project{}, Info: nil, projectsSet: utils.Set[string]{}}
	if _, exists := clients[name]; exists {
		panic("client already exists")
	}
	out.Info = info
	clients[name] = out
	clientsOrder = append(clientsOrder, name)
	return out
}

func (pg *Client) Bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(frontmatterFor(pg.Name, "Client"))
	if pg.Info != nil && len(pg.Info) > 0 {
		builder.WriteString("# Responsibilities and Achievements \n")
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

var (
	jdClient = NewClient("John Deere",
		"Most senior consulting engineer on a high-performance, global-scale, team of 4-8 at John Deere’s Intelligent Solutions Group; responsible for architecture, implementation, testing, optimization, and support of a complex set of diverse cloud services utilizing geospatiotemporal agribusiness data",
		"Architected, implemented, and maintained cloud infrastructure and services for ingest, distributed processing, storage, manipulation, and retrieval of data for, datastores totalling over 50PB",
		"Created multiple ECS clusters for use in, and consuming data from, John Deere AI platforms",
		"Designed CUDA C/C++ Kernels used through CGo for statistics and image processing via GPU",
		"Improved a mission-critical image manipulation API from 65% reliability to 99.9999% success rate",
		"Spearheaded implementation of an ECS cluster using an advanced topology algorithm (from a PhD thesis), achieving >100x performance gains over its original implementation on TB-scale datasets",
		"Built a Go-based compiler that transformed nested JSON instructions into machine-executable operations across EC2 clusters for for retrieving and manipulating geospatial agricultural data",
		"Saved $18M of a $28M budget (64%) in 2024, while still increasing service stability and throughput",
		"Ensured maximum service uptime via careful design and rollout of CI/CD pipelines operated via self-hosted GitHub Actions runners, in conjunction with Infrastructure as Code via Terraform",
		"Setup  monitoring, dashboards, alerting, traces, and profiling (Datadog/Grafana/CloudWatch)",
		"Responsible for educating engineers on infrastructure, codebase, domain, and best practices",
		"Utilized primarily Go, Terraform, Bash, and Docker on AWS, but also used Scala, Github Actions, C/C++ with CUDA, DroneCI, python, javascript, typescript, Kotlin, and more",
		"Platforms utilized:  AWS (>30 separate services), Datadog, LogCentral, Rally, Azure DevOps, Github, DroneCI, Grafana, Prometheus, Confluence, and more",
	)
	simpsonUniversityClient = NewClient("Simpson University",
		"Upgraded the University’s payment and donation gateway. Remediated resulting bugs",
	)
	sourceAlliesClient = NewClient("Source Allies",
		"Source Allies internal projects",
		"Upgraded company internal payment gateway to a newer version of Java Spring"
		"Secured all company machines via JAMF to ensure the protection of company and client data",
		"Designed and created a Slack bot integration with Small Improvements to automate monthly announcements, and employee creation and completion of personal and professional goals",
		"Redesigned Jira workflows streamlining the hiring process and onboarding systems for remote coworkers",
	)
	critColaClient = NewClient("CritCola",
		"Consulted on hosting game servers and discord bots for a large online community such that they could be deployed or destroyed, on short notice with persistent game data utilizing GitLabCI, Terraform, CloudFlare, and AWS (EC2, EBS, IAM)",
		"Created a final product with a spin-up time of approximately 3 minutes",
	)
	arrowNailClient = NewClient("ArrowNail LLC", // TODO: ARROWNAIL CLIENT
		"Used serverless services on AWS to support a React geospatial web app, Node API via lambda functions, and an aurora database. Provided IT and Systems Administration Support",
		"Set up CI/CD pipeline in Gitlab CI to make future deployments seamless",
	)
	wellAwareClient = NewClient("Well Aware NC",
		"Designed, created, and hosted a website for Well Aware NC, a University of North Carolina Chapel Hill affiliated nonprofit focused on the testing of well water contaminants within North Carolina",
	)
	clarkClient = NewClient("Chapel Hill Masters Student in Public Health",
		"Created programs for a student doing research for his Masters Degree in Public Health. Provided data on E.Coli samples from different waterways, the programs checked the statistical validity on different E.Coli indicating kits",
	)
	wildlifeRClient = NewClient("Wildlife mapping with R", // TODO: WILDLIFE MAPPING IN R: .
		"Utilized spatiotemporal data for wildlife in a specified area over a specified date range in order to produce population density maps",
	)
	charityClient = NewClient("Undisclosed Charity",
		"Designed, created, and hosted a website for an undisclosed local charity",
	)
	teiClient = NewClient("TEI",
		"Internal company work for TEI",
	)
	taeClient = NewClient("Talley Associates of Engineering", // TODO: JS photo parser
		"Internal company work for Talley Associates of Engineering",
	)
	mafcClient = NewClient("Monroe Aquatics and Fitness Center",
		"Lifeguarding work, both at the indoor and outdoor pools of the fitness center",
	)
)

func initClientsAfterProjectsComplete() {
	jdClient = jdClient.WithProjects(polygonBuilderProject, ogreProject, renderProject, statsProject, billingProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, tileGenProject, GhaRunnersProject)
	sourceAlliesClient = sourceAlliesClient.WithProjects(simpsonUnivProject) // TODO: USE!
	// TODO: SIMPSON COLLEGE // TODO: USE!
	critColaClient = critColaClient.WithProjects(CritColaProject)
	wellAwareClient = wellAwareClient.WithProjects(WellAwareProject)
	clarkClient = clarkClient.WithProjects(MastersDataAnalysisProject)
	charityClient = charityClient.WithProjects(CharityProject)
	teiClient = teiClient.WithProjects(teiProjects)    // TODO: add projects (like NM, TX, IA, NC?)
	taeClient = taeClient.WithProjects(taeProjects)    // TODO: JS photo parser
	mafcClient = mafcClient.WithProjects(mafcProjects) // TODO: Indoor and outdoor pool?
	// TODO: wildlifeRClient, arrowNailClient, simpsonUniversityClient
}
