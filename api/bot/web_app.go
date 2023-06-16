package bot

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"html"
	"net/http"
	"text/template"
)

var (
	//go:embed index.html
	Resource embed.FS
	//go:embed login.html
	LoginResource embed.FS
	WebApp        = hWebApp{}
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

func (h *hWebApp) CheckAuthorization(r *ghttp.Request) {
	query := r.Request.URL.Query()
	token := "5548720536:AAFY-wb4ir22eF5vRMQXft_sj-RDhaB54EQ"
	ok, err := ValidateLoginQuery(query, token)
	if err != nil {
		return
	}
	if ok {
		saveTelegramUserData(r)
	}
	// 重定向
	r.Response.RedirectTo("/login", http.StatusFound)
}

type TGUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func saveTelegramUserData(r *ghttp.Request) {
	q := r.Request.URL.Query()
	data := TGUser{
		Id:        gconv.Int64(q.Get("id")),
		FirstName: q.Get("first_name"),
		LastName:  q.Get("last_name"),
		Username:  q.Get("username"),
		PhotoUrl:  q.Get("photo_url"),
		AuthDate:  gconv.Int64(q.Get("auth_date")),
		Hash:      q.Get("hash"),
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Println("saveTelegramUserData result:", string(bytes))
	r.Cookie.Set("tg_user", base64.StdEncoding.EncodeToString(bytes))
}

func getTelegramUserData(r *ghttp.Request) (*TGUser, error) {
	ustr := r.Cookie.Get("tg_user")
	if ustr.String() != "" {
		var authData TGUser
		str, err := base64.StdEncoding.DecodeString(ustr.String())
		if err != nil {
			return nil, err
		}
		fmt.Println("getTelegramUserData result:", string(str))
		err = json.Unmarshal(str, &authData)
		if err != nil {
			return nil, err
		}
		return &authData, nil
	}
	return nil, errors.New("not Found")
}

func (h *hWebApp) Login(r *ghttp.Request) {
	if log := r.Get("logout"); log.String() != "" {
		r.Cookie.Set("tg_user", "")
		fmt.Println("logout .....")
		r.Response.RedirectTo("/login", http.StatusMovedPermanently)
		return
	}
	tgUser, err := getTelegramUserData(r)
	if err != nil {
		fmt.Println("getTelegramUserData  error:", err.Error())
		//return
	}
	htm := ""
	if tgUser != nil {
		username := html.EscapeString(tgUser.Username)
		firstName := html.EscapeString(tgUser.FirstName)
		lastName := html.EscapeString(tgUser.LastName)
		photoUrl := html.EscapeString(tgUser.PhotoUrl)
		if tgUser.Username != "" {
			tpl := `<h1>Hello, <a href="https://t.me/%s">%s %s</a>!</h1>`
			htm = fmt.Sprintf(tpl, username, firstName, lastName)
		} else {
			tpl := `<h1>Hello, %s %s!</h1>`
			htm = fmt.Sprintf(tpl, firstName, lastName)
		}
		if tgUser.PhotoUrl != "" {
			tpl := `<img src="%s">`
			htm = fmt.Sprintf("%s\n%s", htm, fmt.Sprintf(tpl, photoUrl))
		}
		tpl := `<p><a href="?logout=1">Log out</a></p>`
		htm = fmt.Sprintf("%s\n%s", htm, tpl)
	} else {
		botUsername := "smartWk_bot"
		htm = `
<h1>Hello, anonymous!</h1>
<script async src="https://telegram.org/js/telegram-widget.js?22" data-telegram-login="%s" data-size="large" data-auth-url="/check_authorization"></script>
`
		htm = fmt.Sprintf(htm, botUsername)
	}

	indexTmpl := template.Must(template.ParseFS(LoginResource, "login.html"))
	writer := r.Response.ResponseWriter
	err = indexTmpl.ExecuteTemplate(writer, "login.html", struct {
		Html string
	}{
		Html: htm,
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
	}
}
