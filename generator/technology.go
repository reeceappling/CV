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
	NewTechnology("RFID").
		WithSubjectMatters(smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("NFC").
		WithSubjectMatters(smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("I2C").
		WithSubjectMatters(smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("SPI").
		WithSubjectMatters(smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("Drupal").
		WithSubjectMatters(smFrontend)
	NewTechnology("React").
		WithSubjectMatters(smFrontend)
	NewTechnology("Quartz 4").
		WithSubjectMatters(smFrontend)
	NewTechnology("Quartz").
		WithSubjectMatters(smFrontend)
	NewTechnology("NextJs").
		WithSubjectMatters(smFullStack)
	NewTechnology("GraphQL").
		WithSubjectMatters(smBackend, smNetworking)
	NewTechnology("Obsidian").
		WithSubjectMatters(smFrontend)
	NewTechnology("SIMD").
		WithSubjectMatters(smBackend, smRobotics, smEmbeddedSystems)
	NewTechnology("Cloudflare Tunnels").
		WithSubjectMatters(smNetworking)
	NewTechnology("LocalStack").
		WithSubjectMatters(smBackend, smCloudComputing, smNetworking)
	NewTechnology("ElasticSearch").
		WithSubjectMatters(smBackend)
	NewTechnology("NodeJS").
		WithSubjectMatters(smBackend)
	NewTechnology("Markdown").
		WithSubjectMatters(smObservability)
	NewTechnology("Websockets").
		WithSubjectMatters(smFullStack, smNetworking)
	NewTechnology("Server-Sent Events").
		WithSubjectMatters(smFullStack, smNetworking)
	NewTechnology("SIMULATE3").
		WithSubjectMatters(smNuclearEngineering, smParticlePhysics)
	NewTechnology("CASMO4e").
		WithSubjectMatters(smNuclearEngineering, smParticlePhysics)
	NewTechnology("Pub-Sub").
		WithSubjectMatters(smBackend, smNetworking)
	NewTechnology("Git").
		WithSubjectMatters(smFullStack, smDevOps)
	NewTechnology("Webhooks").
		WithSubjectMatters(smFullStack, smNetworking)
	NewTechnology("JQuery").
		WithSubjectMatters(smFrontend)
	NewTechnology("Arduino").
		WithSubjectMatters(smEmbeddedSystems, smElectronics, smRobotics)
	NewTechnology("PWM").
		WithSubjectMatters(smEmbeddedSystems, smElectronics, smRobotics)
}

// TODO: list all backlinks????
