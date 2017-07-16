package templates

import (
	"errors"
	"html/template"
	"io"
	"sync"
)

var NotFound error

func init() {
	NotFound = errors.New("Item missing")
	configs = make(map[string]*Config)
	hosts = make(map[string]string)
	master = &sync.RWMutex{}
}

var (
	FirstDelims  = [2]string{"[[", "]]"}
	SecondDelims = [2]string{"{{", "}}"}
)

var (
	configs map[string]*Config
	hosts   map[string]string
	master  *sync.RWMutex
)

type Config struct {
	Name  string
	Hosts map[string]*Host
	Lock  *sync.RWMutex
	TemplateData
}

type TemplateData struct {
	Masters      map[string]*MasterTemplate
	Features     map[string]*FeatureTemplate
	FeatureLists map[string][]string
}

type Host struct {
	Name     string
	Config   map[string]interface{}
	Compiled map[string]*template.Template
	Features map[string]bool
	Lock     *sync.RWMutex
	TemplateData
}

func (c *Config) GenerateFor(n string) error {
	if h, ok := c.Hosts[n]; ok {
		h.Lock.Lock()
		defer h.Lock.Unlock()
		c.Lock.RLock()
		defer c.Lock.RUnlock()
		for _, mt := range c.Masters {
			err := mt.Generate(h)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) Regenerate() error {
	{
		c.Lock.RLock()
		defer c.Lock.RUnlock()
		for h := range c.Hosts {
			err := c.GenerateFor(h)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Host) Render(name string, w io.Writer, args interface{}) error {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	if t, ok := h.Compiled[name]; ok {
		return t.Execute(w, args)
	}
	return NotFound
}
