package bot

import (
	"embed"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"text/template"
)

var (
	//go:embed index.html
	Resource embed.FS
	WebApp   = hWebApp{}
)

type hWebApp struct {
	Path     string `json:"path" yaml:"path"`
	BasePath string `json:"basePath" yaml:"base_path"`
}

func (h *hWebApp) Index(r *ghttp.Request) {
	indexTmpl := template.Must(template.ParseFS(Resource, "index.html"))
	writer := r.Response.ResponseWriter
	err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
		WebAppURL string
	}{
		WebAppURL: "https://dev.us.jxo.me",
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
	}
}

func (h *hWebApp) Validate(r *ghttp.Request) {
	writer := r.Response.ResponseWriter
	token := "5548720536:AAFY-wb4ir22eF5vRMQXft_sj-RDhaB54EQ"
	ok, err := ValidateWebAppQuery(r.Request.URL.Query(), token)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("validation failed; error: " + err.Error()))
		return
	}
	if ok {
		writer.Write([]byte("validation success; user is authenticated."))
	} else {
		writer.Write([]byte("validation failed; data cannot be trusted."))
	}
}
