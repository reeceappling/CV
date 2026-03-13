package main

import "appli.ng/cv/generator/utils"

var techs = map[string]*TechologyPage{}

type TechologyPage struct { // React, Github actions, etc
	Name string
	*tracked
	SubjectMatters utils.Set[SubjectMatter]
}

func (pg *TechologyPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "technology", withoutSpaces(pg.Name))
}

func (pg *TechologyPage) EntryType() string {
	return "Technology"
}

func (pg *TechologyPage) WithSubjectMatters(sms ...SubjectMatter) *TechologyPage {
	if pg == nil {
		return nil
	}
	pg.SubjectMatters = withSubjectMatters(pg, pg.SubjectMatters, sms...)
	return pg
}

func (pg *TechologyPage) WithTags(tags ...string) *TechologyPage {
	if pg == nil {
		return nil
	}
	pg.Tags.Add(tags...) // TODO: USE THESE TAGS
	return pg
}

func NewTechnology(name string) *TechologyPage {
	out := &TechologyPage{
		Name:           name,
		tracked:        newTracked(),
		SubjectMatters: utils.Set[SubjectMatter]{},
	}
	techs[name] = out
	return out
}

func setupTechSubjectMatters() {
	NewTechnology("Github Actions").
		WithSubjectMatters(smCiCd)
	NewTechnology("Gitlab CI").
		WithSubjectMatters(smCiCd)
	NewTechnology("Drone CI").
		WithSubjectMatters(smCiCd)
	NewTechnology("LangGraph").
		WithSubjectMatters(smAI)
	NewTechnology("LangChain").
		WithSubjectMatters(smAI)
	NewTechnology("LLM").
		WithSubjectMatters(smAI)
	NewTechnology("AI").
		WithSubjectMatters(smAI)
	NewTechnology("AI Agents").
		WithSubjectMatters(smAI)
	NewTechnology("OpenAI API spec").
		WithSubjectMatters(smAI, smBackend)
	NewTechnology("CUDA").
		WithSubjectMatters(smGPU, smGraphics)
	NewTechnology("Avro").
		WithSubjectMatters(smBackend).
		WithTags("DataFormat")
	NewTechnology("Parquet").
		WithSubjectMatters(smBackend).
		WithTags("DataFormat")
	NewTechnology("CUDA").
		WithSubjectMatters(smGPU, smGraphics)
	// TODO: GRAPHQL????
	NewTechnology("Kubernetes").
		WithSubjectMatters(smBackend, smDistributedComputing, smContainerization, smIAC, smCloudComputing, smNetworking)
	// TODO: ansible? chef?
}

// TODO: list all backlinks????
