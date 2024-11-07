package go_http_hotreload

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

func AppendToServeMux(mux *http.ServeMux) error {

	mux.HandleFunc("GET /hotreload.js", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := boilerplate.ReadFile("embed/hotreload.enabled.js")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("failed to load hotreload.enabled.js"))
		} else {
			w.Write(bytes)
		}
	})

	mux.HandleFunc("HEAD /__hotreload.ws", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("GET /__hotreload.ws", (websocket.Server{
		Handler: func(c *websocket.Conn) {
			for {
				time.Sleep(1000 * time.Millisecond)
				_, err := c.Write([]byte("Ping"))
				if err != nil {
					return
				}
			}
		},
	}).ServeHTTP)

	return nil

}
