package main

import (
	"log"
	"net/http"

	"hydrz.com/go-vue-embed/dist"
)

func main() {
	http.Handle("/", dist.NewAssetsHandler())

	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		var s = "test"
		rw.Write([]byte(s))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
