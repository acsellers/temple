package templates

import (
	"bytes"
	ttemplate "text/template"
)

type FeatureTemplate struct {
	Name     string
	Template *ttemplate.Template
	Config   *Config
}

func (ft *FeatureTemplate) Generate(h *Host) string {
	args := make(map[string]interface{})
	args["Features"] = func(name string) string {
		result := ""
		for _, t := range ft.Config.FeatureLists[name] {
			if h.Features[t] {
				result += ft.Config.Features[t].Generate(h)
			}
		}
		return result
	}
	b := &bytes.Buffer{}
	ft.Template.Execute(b, args)
	return b.String()
}
