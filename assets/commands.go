package assets

import (
	"fmt"
	"sync"
	"time"
)

func NewConfig(config string) error {
	master.RLock()
	// already exists
	_, ok := configs[config]
	master.RUnlock()
	if ok {
		return nil
	}

	c := &Config{
		Name:         config,
		Hosts:        make(map[string]*Host),
		Lock:         &sync.RWMutex{},
		AssetFolders: make(map[string]*AssetFolder),
	}
	master.Lock()
	// in case of two identical requests
	if _, ok := configs[config]; !ok {
		configs[config] = c
	}
	master.Unlock()

	return nil
}

func DeleteConfig(config string) error {
	master.Lock()
	_, ok := configs[config]
	if ok {
		delete(configs, config)
		master.Unlock()
	} else {
		master.Unlock()
		return NotFound
	}

	return nil
}

func SaveAssetFolder(config, host, folder string, headers map[string]string) error {
	master.RLock()
	c, ok := configs[config]
	master.RUnlock()
	if !ok {
		return NotFound
	}
	if host == "" {
		c.Lock.Lock()
		f, ok := c.AssetFolders[folder]
		if ok {
			if d, ok := headers["Expires"]; ok {
				bd, err := NewBigDuration(d)
				if err == nil {
					f.Expires = bd
					delete(headers, "Expires")
				}
			}
			f.Headers = headers
		} else {
			a := &AssetFolder{
				Assets:  make(map[string][]byte),
				ETags:   make(map[string]string),
				Headers: make(map[string]string),
			}
			if d, ok := headers["Expires"]; ok {
				bd, err := NewBigDuration(d)
				if err == nil {
					a.Expires = bd
					delete(headers, "Expires")
				}
			}
			c.AssetFolders[folder] = f
		}
		c.Lock.Unlock()
	} else {
		c.Lock.Lock()
		h, ok := c.Hosts[host]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		f, ok := h.Overrides[folder]
		if ok {
			if d, ok := headers["Expires"]; ok {
				bd, err := NewBigDuration(d)
				if err == nil {
					f.Expires = bd
					delete(headers, "Expires")
				}
			}
			f.Headers = headers
		} else {
			a := &AssetFolder{
				Assets:  make(map[string][]byte),
				ETags:   make(map[string]string),
				Headers: make(map[string]string),
			}
			if d, ok := headers["Expires"]; ok {
				bd, err := NewBigDuration(d)
				if err == nil {
					a.Expires = bd
					delete(headers, "Expires")
				}
			}
			h.Overrides[folder] = a
		}
		c.Lock.Unlock()
	}
	return nil
}

func DeleteAssetFolder(config, host, folder string) error {
	master.RLock()
	c, ok := configs[config]
	master.RUnlock()
	if !ok {
		return NotFound
	}
	if host == "" {
		c.Lock.Lock()
		_, ok := c.AssetFolders[folder]
		if ok {
			delete(c.AssetFolders, folder)
			c.Lock.Unlock()
			return nil
		} else {
			c.Lock.Unlock()
			return NotFound
		}
	} else {
		c.Lock.Lock()
		h, ok := c.Hosts[host]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		_, ok = h.Overrides[folder]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		delete(h.Overrides, folder)
		c.Lock.Unlock()
	}

	return nil
}

func SaveAsset(config, host, folder, name string, data []byte) error {
	master.RLock()
	c, ok := configs[config]
	master.RUnlock()
	if !ok {
		return NotFound
	}
	if host == "" {
		c.Lock.Lock()
		f, ok := c.AssetFolders[folder]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		f.Assets[name] = data
		f.ETags[name] = fmt.Sprintf("%x", time.Now().Unix())
		c.Lock.Unlock()
	} else {
		c.Lock.Lock()
		h := c.Hosts[host]
		f, ok := h.Overrides[folder]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		f.Assets[name] = data
		f.ETags[name] = fmt.Sprintf("%x", time.Now().Unix())
		c.Lock.Unlock()
	}
	return nil
}

func DeleteAsset(config, host, folder, name string) error {
	master.RLock()
	c, ok := configs[config]
	if !ok {
		master.RUnlock()
		return nil
	}
	master.RUnlock()

	c.Lock.Lock()
	if host == "" {
		f, ok := c.AssetFolders[folder]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		if _, ok := f.Assets[name]; !ok {
			c.Lock.Unlock()
			return NotFound
		}
		delete(f.Assets, name)
		delete(f.ETags, name)
	} else {
		h, ok := c.Hosts[host]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		f, ok := h.Overrides[folder]
		if !ok {
			c.Lock.Unlock()
			return NotFound
		}
		if _, ok := f.Assets[name]; !ok {
			c.Lock.Unlock()
			return NotFound
		}
		delete(f.Assets, name)
		delete(f.ETags, name)
	}
	c.Lock.Unlock()

	return nil
}
