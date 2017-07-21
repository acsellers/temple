package main

import (
	"bytes"
	htemplate "html/template"
	"sort"
	ttemplate "text/template"
)

type MasterTemplate struct {
	Template *ttemplate.Template
	Features map[string][]FeatureTemplate
}

func (mt *MasterTemplate) Generate(c *Client) string {
	sort.Strings(c.Features)
	var err error
	site := htemplate.New("layout")
	args := make(map[string]interface{})
	args["Features"] = func(name string) string {
		result := ""
		for _, t := range mt.Features[name] {
			result += t.Generate(c)
		}
		return result
	}
	b := &bytes.Buffer{}
	mt.Template.Execute(b, args)
	c.Site, err = site.Parse(b.String())
	if err != nil {
		panic(err)
	}
	return b.String()
}

type FeatureTemplate struct {
	Master   *MasterTemplate
	Name     string
	Template *ttemplate.Template
}

func (ft *FeatureTemplate) Generate(c *Client) string {
	if len(c.Features) == sort.SearchStrings(c.Features, ft.Name) {
		return ""
	}
	args := make(map[string]interface{})
	args["Features"] = func(name string) string {
		result := ""
		for _, t := range ft.Master.Features[name] {
			result += t.Generate(c)
		}
		return result
	}
	b := &bytes.Buffer{}
	ft.Template.Execute(b, args)
	return b.String()
}
