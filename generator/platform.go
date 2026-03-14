package main

var platforms = map[string]*PlatformPage{}

type PlatformPage struct { // Datadog, Github, etc
	Name string
	*tracked
	SubjectMattersField
}

func (pg *PlatformPage) EntryType() string {
	return "Platform"
}

func (pg *PlatformPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "platform", withoutSpaces(pg.Name))
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
	NewPlatform("Azure DevOps"). // TODO: USE
		WithSubjectMatters(smDevOps)
	NewPlatform("Jira"). // TODO: USE
		WithSubjectMatters(smDevOps)
	NewPlatform("Confluence"). // TODO: USE
		WithSubjectMatters(smDocumentation)
	NewPlatform("Slack"). // TODO: USE
		WithSubjectMatters(smCommunication)
	NewPlatform("Small Improvements"). // TODO: USE
		WithSubjectMatters(smObservability)
	NewPlatform("Teams").
		WithSubjectMatters(smCommunication)
}
