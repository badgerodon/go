package web

import (
	"net/http"
)

func Listen(addr string) {
	http.Handle(
		"/assets/",
		http.StripPrefix(
			"/assets/",
			http.FileServer(http.Dir("./assets")),
		),
	)
	for route, application := range Applications {
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			application.Handle(w, r)
		})
	}
	http.ListenAndServe(addr, nil)
}