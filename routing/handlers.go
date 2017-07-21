package routing

import (
	"io"
	"net/http"
)

func MainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, err := GetHost(r.URL.Host)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Unrecognized Host")
			return
		}
		route, err := h.GetRoute()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			h.WriteError(404, w, r)
			return
		}
		route.ServeHTTP(w, r)
	})
}

func APIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* API
		/routing/:config/
		/routing/:config/
		/routing/:config/
		/routing/:config/
		POST /assets/:config
		DELETE /assets/:config
		POST /assets/:config/:folder
		DELETE /assets/:config/:folder
		POST /assets/:config/:folder/:asset
		DELETE /assets/:config/:folder/:asset
		/templates/
		*/
	})
}
