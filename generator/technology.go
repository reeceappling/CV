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

func NewTechnology(name string, subjectMatters ...SubjectMatter) *TechologyPage {
	out := &TechologyPage{
		Name:           name,
		tracked:        newTracked(),
		SubjectMatters: utils.SetFrom(subjectMatters...),
	}
	techs[name] = out
	return out
}

func setupTechSubjectMatters() {
	NewTechnology("Github Actions", smCiCd)
	NewTechnology("Gitlab CI", smCiCd)
	NewTechnology("Drone CI", smCiCd)
	NewTechnology("LangGraph", smAI)
	NewTechnology("LangChain", smAI)
	NewTechnology("LLM", smAI)
	NewTechnology("AI", smAI)
	NewTechnology("AI Agents", smAI)
	NewTechnology("OpenAI API spec", smAI, smBackend)
	NewTechnology("CUDA", smGPU, smGraphics)
	NewTechnology("Avro", smBackend).WithTags("DataFormat")
	NewTechnology("Parquet", smBackend).WithTags("DataFormat")
	NewTechnology("JSON", smFullStack).WithTags("DataFormat")
	NewTechnology("ENDF", smNuclearEngineering).WithTags("DataFormat")
	NewTechnology("JEFF", smNuclearEngineering).WithTags("DataFormat")
	NewTechnology("XML", smFullStack).WithTags("DataFormat")
	NewTechnology("YAML", smFullStack).WithTags("DataFormat")
	NewTechnology("TOML", smFullStack).WithTags("DataFormat")
	NewTechnology("CUDA", smGPU, smGraphics)
	// TODO: GRAPHQL????
	NewTechnology("Kubernetes", smBackend, smDistributedComputing, smContainerization, smIAC, smCloudComputing, smNetworking)
	// TODO: ansible? chef?
	NewTechnology("RFID", smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("NFC", smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("I2C", smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("SPI", smRobotics, smEmbeddedSystems, smElectronics)
	NewTechnology("Drupal", smFrontend)
	NewTechnology("React", smFrontend)
	NewTechnology("Quartz 4", smFrontend)
	NewTechnology("Quartz", smFrontend)
	NewTechnology("NextJs", smFullStack)
	NewTechnology("GraphQL", smBackend, smNetworking)
	NewTechnology("Obsidian", smFrontend)
	NewTechnology("SIMD", smBackend, smRobotics, smEmbeddedSystems)
	NewTechnology("Cloudflare Tunnels", smNetworking)
	NewTechnology("LocalStack", smBackend, smCloudComputing, smNetworking)
	NewTechnology("ElasticSearch", smBackend)
	NewTechnology("NodeJS", smBackend)
	NewTechnology("Markdown", smObservability)
	NewTechnology("Websockets", smFullStack, smNetworking)
	NewTechnology("Server-Sent Events", smFullStack, smNetworking)
	// TODO: ADD NDSF(?) FILES FOR NUC STUFF
	NewTechnology("SIMULATE3", smNuclearEngineering, smParticlePhysics)
	NewTechnology("CASMO4e", smNuclearEngineering, smParticlePhysics)
	NewTechnology("Pub-Sub", smBackend, smNetworking)
	NewTechnology("Git", smFullStack, smDevOps)
	NewTechnology("Webhooks", smFullStack, smNetworking)
	NewTechnology("JQuery", smFrontend)
	NewTechnology("Arduino", smEmbeddedSystems, smElectronics, smRobotics)
	NewTechnology("PWM", smEmbeddedSystems, smElectronics, smRobotics)
	NewTechnology("G and M codes", smElectronics, smRobotics)
}

// TODO: list all backlinks????
