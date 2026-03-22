package main

import (
	"appli.ng/cv/generator/utils"
	"strings"
)

var techs = map[string]*TechologyPage{}

type TechologyPage struct { // React, Github actions, etc
	Name string
	*tracked
	SubjectMatters utils.Set[SubjectMatter]
}

func (pg *TechologyPage) Dst() string {
	return dstFor("cv", "technology", withoutSpaces(pg.Name))
}

func (pg *TechologyPage) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
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
	addLinkable(strings.ToLower(name), out)
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
	NewTechnology("OpenAI API spec", smAI, smBackend).WithTags("API")
	NewTechnology("CUDA", smGPU, smGraphics)
	NewTechnology("CGo", smBackend)
	NewTechnology("REST API", smBackend, smApi).WithTags("API") // TODO: use this everywhere necessary...
	NewTechnology("Avro", smBackend).WithTags("DataFormat")
	NewTechnology("Parquet", smBackend).WithTags("DataFormat")
	NewTechnology("JSON", smFullStack).WithTags("DataFormat")
	NewTechnology("ENDF", smNuclearEngineering).WithTags("DataFormat")
	NewTechnology("JEFF", smNuclearEngineering).WithTags("DataFormat")
	NewTechnology("XML", smFullStack).WithTags("DataFormat")
	NewTechnology("YAML", smFullStack, smCiCd).WithTags("DataFormat")
	NewTechnology("TOML", smFullStack).WithTags("DataFormat")
	NewTechnology("CUDA", smGPU, smGraphics)
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
	NewTechnology("GraphQL", smBackend, smNetworking).WithTags("API")
	NewTechnology("Obsidian", smDocumentation)
	NewTechnology("SIMD", smBackend, smRobotics, smEmbeddedSystems)
	NewTechnology("Cloudflare Tunnels", smNetworking)
	NewTechnology("LocalStack", smBackend, smCloudComputing, smNetworking)
	NewTechnology("ElasticSearch", smBackend)
	NewTechnology("NodeJS", smBackend)
	NewTechnology("Markdown", smDocumentation)
	NewTechnology("Websockets", smFullStack, smNetworking)
	NewTechnology("Server-Sent Events", smFullStack, smNetworking)
	NewTechnology("SIMULATE3", smNuclearEngineering, smParticlePhysics)
	NewTechnology("CASMO4e", smNuclearEngineering, smParticlePhysics)
	NewTechnology("Pub-Sub", smBackend, smNetworking)
	NewTechnology("Git", smFullStack, smDevSecOps)
	NewTechnology("Webhooks", smFullStack, smNetworking)
	NewTechnology("JQuery", smFrontend)
	NewTechnology("Arduino", smEmbeddedSystems, smElectronics, smRobotics)
	NewTechnology("PWM", smEmbeddedSystems, smElectronics, smRobotics)
	NewTechnology("G and M codes", smElectronics, smRobotics)
	NewTechnology("OpenApi", smDocumentation).WithTags("API")
	NewTechnology("Swagger", smDocumentation).WithTags("API")
	NewTechnology("Spring", smBackend)
	//NewTechnology("Spark", smBackend) // TODO: apache arrow
	//NewTechnology("Arrow", smBackend) // TODO: apache arrow
	//NewTechnology("ORC", smBackend) // TODO: apache ORC
	NewTechnology("Kafka", smBackend) // TODO: event driven?
	//
}

// TODO: list all backlinks????
