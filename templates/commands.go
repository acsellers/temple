package templates

import (
	"sync"
	ttemplate "text/template"
)

func SaveMaster(config, name, tmpl string) error {
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
		{
			c.Lock.Lock()
			c.Masters[name] = nm
			c.Lock.Unlock()
		}
		return c.Regenerate()
	}
	return nil
}

func RemoveMaster(config, name string) error {
	if c, ok := configs[config]; ok {
		c.Lock.Lock()
		delete(c.Masters, name)
		c.Lock.Unlock()
		for _, h := range c.Hosts {
			h.Lock.Lock()
			delete(h.Compiled, name)
			h.Lock.Unlock()
		}
		return nil
	}
	return NotFound
}

func SaveFeature(config, name, tmpl string) error {
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
		{
			c.Lock.Lock()
			c.Features[name] = nf
			c.Lock.Unlock()
		}
		return c.Regenerate()
	}
	return nil
}

func DeleteFeature(config, name string) error {
	if c, ok := configs[config]; ok {
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
	}
	return NotFound
}

func SaveHost(config, host string, hconfig map[string]interface{}, features []string) error {
	if c, ok := configs[config]; ok {
		h := &Host{Config: hconfig, Name: host, Features: make(map[string]bool), Lock: &sync.RWMutex{}}
		for _, f := range features {
			h.Features[f] = true
		}
		// generate compiled templates
		// need to lock compiled templates to ensure that
		// we're not going to be outdated when we insert
		// the new host
		c.Lock.RLock()
		defer c.Lock.RUnlock()
		for _, mt := range c.Masters {
			err := mt.Generate(h)
			if err != nil {
				return err
			}
		}
		{ // replace host entry
			c.Lock.Lock()
			c.Hosts[host] = h
			c.Lock.Unlock()
		}
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
			Name:         config,
			Hosts:        make(map[string]*Host),
			Lock:         &sync.RWMutex{},
			Masters:      make(map[string]*MasterTemplate),
			Features:     make(map[string]*FeatureTemplate),
			FeatureLists: make(map[string][]string),
		}
		configs[config] = c
	}

	return nil
}

func DeleteConfig(config string) error {
	if _, ok := configs[config]; ok {
		delete(configs, config)
	}
	return nil
}
