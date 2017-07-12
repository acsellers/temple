package main

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"log"
	ttemplate "text/template"
)

type Client struct {
	Features []string
	Name     string
	Site     *htemplate.Template
}

var MT *MasterTemplate

func init() {
	MT = &MasterTemplate{}
	var err error
	var mt *ttemplate.Template
	mt, err = ttemplate.New("master").Delims("[[", "]]").Parse(Master)
	if err != nil {
		log.Fatal(err)
	}
	MT.Template = mt
	MT.Features = make(map[string][]FeatureTemplate)
	for name, top := range Tops {
		tt, err := ttemplate.New(name).Parse(top)
		if err != nil {
			log.Fatal(name, err)
		}
		MT.Features["topbar"] = append(
			MT.Features["topbar"],
			FeatureTemplate{Name: name, Template: tt},
		)
	}
	for name, nav := range Navs {
		tt, err := ttemplate.New(name).Parse(nav)
		if err != nil {
			log.Fatal(name, err)
		}
		MT.Features["navbar"] = append(
			MT.Features["navbar"],
			FeatureTemplate{Name: name, Template: tt},
		)
	}
	for name, foot := range Foots {
		tt, err := ttemplate.New(name).Parse(foot)
		if err != nil {
			log.Fatal(name, err)
		}
		MT.Features["footer"] = append(
			MT.Features["footer"],
			FeatureTemplate{Master: MT, Name: name, Template: tt},
		)
	}
}

func main() {
	c := &Client{Features: []string{"navbar.one", "topbar.three", "footer.two"}}
	fmt.Println(MT.Generate(c))
	b := &bytes.Buffer{}
	c.Site.Execute(b, map[string]interface{}{
		"Title":   "Test",
		"Content": "<script src='/test.js'></script>",
	})
	fmt.Println(b.String())
}
