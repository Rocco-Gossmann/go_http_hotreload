package go_http_hotreload

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func AppendToServeMux(mux *http.ServeMux) error {

	idPool := 0

	mux.HandleFunc("GET /hotreload.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/javascript")
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

			idPool++
			var id = idPool

			// Keep connection awake
			go func() {
				for {
					time.Sleep(1000 * time.Millisecond)
					_, err := c.Write([]byte("Ping"))
					if err != nil {
						return
					}
				}
			}()

			buff := make([]byte, 1024)
			for {
				_, err := c.Read(buff)
				if err == nil {
					log.Println("conn ", id, " send: ", string(buff))
				}
			}

		},
	}).ServeHTTP)

	return nil

}
