package routing

import (
	"errors"
	"net/http"
)

type Config struct {
	Name  string
	Hosts map[string]*Host
}

type Host struct {
	Name   string
	Config *Config
}

func (h *Host) GetRoute() (Route, error) {
	return nil, errors.New("TODO")
}

func (h *Host) WriteError(code int, w http.ResponseWriter, r *http.Request) {
}

type Route interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func GetHost(h string) (*Host, error) {
	return nil, errors.New("TODO")
}

type AssetRoute struct {
}

func (ar AssetRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Lookup Asset
	// Return Asset
}

type ProxyRoute struct {
}

func (ar ProxyRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Proxy Request
}

type TemplateRoute struct {
}

func (ar TemplateRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Proxy Request
	// Check Result
	// Return if not json
	// Find Template
}

type AuthRoute struct {
}

func (ar AuthRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Proxy Things
}
