package templates

import (
	htemplate "html/template"
	"io"
	"sync"
	ttemplate "text/template"
)

func SaveMaster(config, host, name, tmpl string) error {
	if c, ok := configs[config]; ok {
		nm := &MasterTemplate{Name: name, Config: c}
		nt, err := ttemplate.
			New(name).
			Delims(FirstDelims[0], FirstDelims[1]).
			Parse(tmpl)
		if err != nil {
			return err
		}
		nm.Template = nt
		if host == "" {
			c.Lock.Lock()
			c.Masters[name] = nm
			c.Lock.Unlock()
		} else {
			for _, h := range c.Hosts {
				if h.Name == host {
					h.Masters[name] = nm
				}
			}
		}
		return c.Regenerate()
	}
	return nil
}

func RemoveMaster(config, host, name string) error {
	if c, ok := configs[config]; ok {
		c.Lock.Lock()
		if host == "" {
			delete(c.Masters, name)
		} else {
			for _, h := range c.Hosts {
				if host == h.Name {
					h.Lock.Lock()
					delete(h.Masters, name)
					delete(h.Compiled, name)
					h.Lock.Unlock()
				}
			}
		}
		c.Lock.Unlock()
		if host == "" {
			for _, h := range c.Hosts {
				h.Lock.Lock()
				delete(h.Compiled, name)
				h.Lock.Unlock()
			}
		}
		return nil
	}
	return NotFound
}

func SaveFeature(config, host, name, tmpl string) error {
	if c, ok := configs[config]; ok {
		nf := &FeatureTemplate{Name: name, Config: c}
		nt, err := ttemplate.
			New(name).
			Delims(FirstDelims[0], FirstDelims[1]).
			Parse(tmpl)
		if err != nil {
			return err
		}
		nf.Template = nt
		if host == "" {
			c.Lock.Lock()
			c.Features[name] = nf
			c.Lock.Unlock()
			return c.Regenerate()
		} else {
			c.Lock.RLock()
			for _, h := range c.Hosts {
				if h.Name == host {
					h.Lock.Lock()
					h.TemplateData.Features[name] = nf
					h.Lock.Unlock()
					return c.GenerateFor(h.Name)
				}
			}
			c.Lock.RUnlock()
		}
	}
	return nil
}

func DeleteFeature(config, host, name string) error {
	if c, ok := configs[config]; ok {
		if host == "" {
			c.Lock.Lock()
			delete(c.Features, name)
			for n, l := range c.FeatureLists {
				nl := make([]string, 0, len(l)-1)
				for _, li := range l {
					if li != name {
						nl = append(nl, li)
					}
				}
				c.FeatureLists[n] = nl
			}
			c.Lock.Unlock()
			return c.Regenerate()
		} else {
			c.Lock.RLock()
			_, override := c.Features[name]
			if override {
			} else {
			}
		}
	}
	return NotFound
}

func SaveHost(config, host string, hconfig map[string]interface{}, features []string) error {
	if c, ok := configs[config]; ok {
		h := &Host{
			Name:     host,
			Config:   hconfig,
			Features: make(map[string]bool),
			Lock:     &sync.RWMutex{},
			Compiled: make(map[string]*htemplate.Template),
			TemplateData: TemplateData{
				Masters:      make(map[string]*MasterTemplate),
				Features:     make(map[string]*FeatureTemplate),
				FeatureLists: make(map[string][]string),
			},
		}
		for _, f := range features {
			h.Features[f] = true
		}
		// generate compiled templates
		// need to lock compiled templates to ensure that
		// we're not going to be outdated when we insert
		// the new host
		c.Lock.RLock()
		for _, mt := range c.Masters {
			err := mt.Generate(h)
			if err != nil {
				c.Lock.RUnlock()
				return err
			}
		}
		c.Lock.RUnlock()
		c.Lock.Lock()
		c.Hosts[host] = h
		c.Lock.Unlock()
		master.Lock()
		hosts[host] = config
		master.Unlock()
		return nil
	}
	return NotFound
}

func DeleteHost(config, host string) error {
	if c, ok := configs[config]; ok {
		if h, ok := c.Hosts[host]; ok {
			h.Lock.Lock()
			delete(c.Hosts, host)
			h.Lock.Unlock()
		}
	}
	return NotFound
}

func CreateConfig(config string) error {
	if _, ok := configs[config]; !ok {
		c := &Config{
			Name:  config,
			Hosts: make(map[string]*Host),
			Lock:  &sync.RWMutex{},
			TemplateData: TemplateData{
				Masters:      make(map[string]*MasterTemplate),
				Features:     make(map[string]*FeatureTemplate),
				FeatureLists: make(map[string][]string),
			},
		}
		master.Lock()
		configs[config] = c
		master.Unlock()
	}

	return nil
}

func DeleteConfig(config string) error {
	if _, ok := configs[config]; ok {
		master.Lock()
		delete(configs, config)
		master.Unlock()
		return nil
	}
	return NotFound
}

func Execute(host, layout string, w io.Writer, args interface{}) error {
	master.RLock()
	cn, ok := hosts[host]
	if !ok {
		return NotFound
	}
	c := configs[cn]
	h := c.Hosts[host]
	h.Lock.RLock()
	master.RUnlock()
	var err error
	if t, ok := h.Compiled[layout]; ok {
		err = t.Execute(w, args)
	} else {
		err = NotFound
	}
	h.Lock.RUnlock()
	return err
}
