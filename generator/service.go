package main

import "appli.ng/cv/generator/utils"

var cloudServices = map[string]*CloudServicePage{}

type CloudServicePage struct {
	Name     string
	provider string // TODO: add provider to page????
	*tracked
	SubjectMattersField
}

func (pg *CloudServicePage) EntryType() string {
	return "Cloud Service"
}

func (pg *CloudServicePage) WithSubjectMatters(sms ...SubjectMatter) *CloudServicePage {
	if pg == nil {
		return nil
	}
	pg.SubjectMatters = withSubjectMatters(pg, pg.SubjectMatters, sms...)
	return pg
}

func (pg *CloudServicePage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "service", withoutSpaces(pg.Name))
}

func NewService(name string, provider string) *CloudServicePage {
	out := &CloudServicePage{
		Name:     name,
		provider: provider,
		tracked:  newTracked(),
		SubjectMattersField: SubjectMattersField{
			SubjectMatters: utils.SetFrom(smCloudComputing),
		},
	}
	subjectMatters[smCloudComputing][out.EntryType()].Add() // TODO: ok?
	cloudServices[name] = out
	return out
}

func setupServiceSubjectMatters() {
	NewService("Cloudwatch", "AWS").
		WithSubjectMatters(smObservability)
	NewService("Lambda", "AWS").
		WithSubjectMatters(smServerless)
	NewService("Fargate", "AWS").
		WithSubjectMatters(smServerless)
	NewService("Aurora", "AWS").
		WithSubjectMatters(smServerless)
	NewService("DynamoDB", "AWS").
		WithSubjectMatters(smServerless)
	NewService("SAM", "AWS"). // TODO: use somewhere
					WithSubjectMatters(smServerless)
	NewService("Cloudformation", "AWS"). // TODO: use somewhere
						WithSubjectMatters(smIAC)

}
