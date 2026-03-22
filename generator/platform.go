package main

import "strings"

var platforms = map[string]*PlatformPage{}

type PlatformPage struct { // Datadog, Github, etc
	Name string
	*tracked
	SubjectMattersField
}

func (pg *PlatformPage) EntryType() string {
	return "Platform"
}

func (pg *PlatformPage) Dst() string {
	return dstFor("cv", "platform", withoutSpaces(pg.Name))
}

func (pg *PlatformPage) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
}

func (pg *PlatformPage) WithSubjectMatters(sms ...SubjectMatter) *PlatformPage {
	if pg == nil {
		return nil
	}
	pg.SubjectMatters = withSubjectMatters(pg, pg.SubjectMatters, sms...)
	return pg
}

func NewPlatform(name string) *PlatformPage {
	out := &PlatformPage{
		Name:                name,
		tracked:             newTracked(),
		SubjectMattersField: SubjectMattersField{SubjectMatters: map[SubjectMatter]struct{}{}},
	}
	platforms[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}

func setupPlatformSubjectMatters() {
	NewPlatform("Poolside AI").
		WithSubjectMatters(smAI)
	NewPlatform("Datadog").
		WithSubjectMatters(smObservability)
	NewPlatform("Logcentral").
		WithSubjectMatters(smObservability, smLogging)
	NewPlatform("Github").
		WithSubjectMatters(smCiCd)
	NewPlatform("ServiceNow").
		WithSubjectMatters(smObservability)
	NewPlatform("Grafana").
		WithSubjectMatters(smObservability)
	NewPlatform("Prometheus").
		WithSubjectMatters(smObservability)
	NewPlatform("GraphQL Apollo").
		WithSubjectMatters(smBackend, smNetworking)
	NewPlatform("Cloudflare").
		WithSubjectMatters(smBackend, smNetworking)
	NewPlatform("Gitlab").
		WithSubjectMatters(smCiCd)
	NewPlatform("Azure DevOps").
		WithSubjectMatters(smDevSecOps)
	NewPlatform("Jira").
		WithSubjectMatters(smDevSecOps)
	NewPlatform("Confluence").
		WithSubjectMatters(smDocumentation)
	NewPlatform("Slack").
		WithSubjectMatters(smCommunication)
	NewPlatform("Small Improvements").
		WithSubjectMatters(smObservability)
	NewPlatform("Teams").
		WithSubjectMatters(smCommunication)
	NewPlatform("Jamf").
		WithSubjectMatters(smCybersecurity)
	NewPlatform("DroneCI").
		WithSubjectMatters(smCiCd)
	NewPlatform("Discord").
		WithSubjectMatters(smCommunication)
	NewPlatform("Rally").
		WithSubjectMatters(smObservability)
}
