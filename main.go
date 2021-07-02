package main

import (
	"log"
	"net/http"

	"hydrz.com/go-vue-embed/assets"
)

func main() {
	http.Handle(assets.AssetsPrefix, http.StripPrefix(assets.AssetsPrefix, http.FileServer(http.FS(assets.AssetsFs))))

	a, err := assets.GetAssets()

	if err != nil {
		panic(err)
	}

	t, err := assets.GetIndexTemplate()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, a)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
