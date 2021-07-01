package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	jsoniter "github.com/json-iterator/go"
	"hydrz.com/embed/assets"
)

//go:embed index.tmpl manifest.json
var index embed.FS

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.AssetsFs))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(index, "index.tmpl")
		if err != nil {
			panic(err)
		}

		var mainJs = "http://localhost:3000/src/main.js"
		var viteJs = "http://localhost:3000/@vite/client"
		var css []string
		var imports []string
		if !strings.Contains(os.Args[0], "/tmp/go-build") {
			viteJs = ""
			manifest, err := index.ReadFile("manifest.json")
			if err != nil {
				panic(err)
			}

			var main = jsoniter.Get(manifest, "src/main.js")
			mainJs = main.Get("file").ToString()

			if main.Get("css").Size() > 0 {
				for i := 0; i < main.Get("css").Size(); i++ {
					css = append(css, main.Get("css", i).ToString())
				}
			}

			if main.Get("imports").Size() > 0 {
				for i := 0; i < main.Get("imports").Size(); i++ {
					var s = main.Get("imports", i).ToString()
					imports = append(imports, jsoniter.Get(manifest, s, "file").ToString())
				}
			}
		}

		t.Execute(w, struct {
			Main    string
			Vite    string
			Css     []string
			Imports []string
		}{
			Main:    mainJs,
			Vite:    viteJs,
			Css:     css,
			Imports: imports,
		})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
