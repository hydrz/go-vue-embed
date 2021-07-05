package dist

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	jsoniter "github.com/json-iterator/go"
)

//go:embed assets/* index.tmpl manifest.json
var assetsFs embed.FS

type assets struct {
	Main string
	Vite string
	Css  []string
}

type assetsHandler struct {
	mode     string
	devHost  string
	assets   assets
	template *template.Template
}

var mode = flag.String("devMode", "prod", "deployment environment")
var devHost = flag.String("devHost", "http://localhost:3000", "vue dev server")

func NewAssetsHandler() *assetsHandler {
	flag.Parse()
	ah := &assetsHandler{}
	ah.mode = *mode
	ah.devHost = *devHost
	ah.getAssets()
	ah.getIndexTemplate()
	return ah
}

// ServeHTTP implements http.Handler.
func (ah *assetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	name := path.Clean(upath)

	var fs http.FileSystem
	if ah.idDev() {
		fs = http.Dir(".")
	} else {
		fs = http.FS(assetsFs)
	}

	f, err := fs.Open(name)

	if err != nil {
		ah.template.Execute(w, ah.assets)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if d.IsDir() {
		ah.template.Execute(w, ah.assets)
		return
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func (ah *assetsHandler) getAssets() {
	var mainJs = ah.devHost + "/vue/main.js"
	var viteJs = ah.devHost + "/@vite/client"
	var css []string
	if !ah.idDev() {
		viteJs = ""
		manifest, err := assetsFs.ReadFile("manifest.json")
		if err != nil {
			panic(err)
		}

		var main = jsoniter.Get(manifest, "vue/main.js")
		mainJs = main.Get("file").ToString()

		if main.Get("css").Size() > 0 {
			for i := 0; i < main.Get("css").Size(); i++ {
				css = append(css, main.Get("css", i).ToString())
			}
		}
	}

	ah.assets = assets{
		Main: mainJs,
		Vite: viteJs,
		Css:  css,
	}
}

func (ah *assetsHandler) getIndexTemplate() {
	t, err := template.ParseFS(assetsFs, "index.tmpl")

	if err != nil {
		panic(err)
	}

	ah.template = t
}

func (ah *assetsHandler) idDev() bool {
	return ah.mode == "dev"
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
