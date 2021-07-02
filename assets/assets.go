package assets

import (
	"embed"
	"os"
	"strings"
	"text/template"

	jsoniter "github.com/json-iterator/go"
)

const AssetsPrefix = "/assets/"

//go:embed *
var AssetsFs embed.FS

type Assets struct {
	Main    string
	Vite    string
	Css     []string
	Imports []string
}

func GetAssets() (a Assets, err error) {
	var mainJs = "http://localhost:3000/vue/main.js"
	var viteJs = "http://localhost:3000/@vite/client"
	var css []string
	var imports []string
	if !strings.Contains(os.Args[0], "/tmp/go-build") {
		viteJs = ""
		manifest, err := AssetsFs.ReadFile("manifest.json")
		if err != nil {
			return a, err
		}

		var main = jsoniter.Get(manifest, "vue/main.js")
		mainJs = AssetsPrefix + main.Get("file").ToString()

		if main.Get("css").Size() > 0 {
			for i := 0; i < main.Get("css").Size(); i++ {
				css = append(css, AssetsPrefix+main.Get("css", i).ToString())
			}
		}

		if main.Get("imports").Size() > 0 {
			for i := 0; i < main.Get("imports").Size(); i++ {
				var s = main.Get("imports", i).ToString()
				imports = append(imports, AssetsPrefix+jsoniter.Get(manifest, s, "file").ToString())
			}
		}
	}

	a = Assets{
		Main:    mainJs,
		Vite:    viteJs,
		Css:     css,
		Imports: imports,
	}

	return a, err
}

func GetIndexTemplate() (t *template.Template, err error) {
	t, err = template.ParseFS(AssetsFs, "index.tmpl")
	return t, err
}
