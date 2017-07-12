package templates

import (
	"bytes"
	htemplate "html/template"
	ttemplate "text/template"
)

type MasterTemplate struct {
	Name     string
	Template *ttemplate.Template
	Config   *Config
}

func (mt *MasterTemplate) Generate(h *Host) error {
	var err error
	site := htemplate.New("layout")
	args := make(map[string]interface{})
	args["Features"] = func(name string) string {
		result := ""
		for _, t := range mt.Config.FeatureLists[name] {
			if h.Features[t] {
				result += mt.Config.Features[t].Generate(h)
			}
		}
		return result
	}
	b := &bytes.Buffer{}
	mt.Template.Execute(b, args)
	h.Compiled[mt.Name], err = site.Parse(b.String())
	if err != nil {
		return err
	}
	return nil
}
