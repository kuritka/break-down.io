package portal

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kuritka/break-down.io/common/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	defaultConfigFile = "config.json"
	defaultLayout = "templates/layout.html"
	templateDir   = "templates/"
)

const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

type IdpConfig struct {
	ClientSecret string `json:"clientSecret"`
	ClientID     string `json:"clientID"`
	CookieStoreKey string `json:"cookieStoreKey"`
	Idp Idp
}

type (
	Idp int
)

const (
	GitHubProvider    Idp = iota
	GoogleHubProvider Idp = iota
)

type Server struct {
	db        *db.CalendarProvider
	router    *mux.Router
	oauthCfg  *oauth2.Config
	store     *sessions.CookieStore
	templates map[string]*template.Template
}

func NewServer(options db.ClientOptions, mux *mux.Router, config *IdpConfig, oauthCfg *oauth2.Config ) *Server {
	db := db.NewDb(options)
	cookieStore := sessions.NewCookieStore([]byte(config.CookieStoreKey))
	templates := map[string]*template.Template{}
	templates["home.html"] = template.Must(template.ParseFiles(templateDir+"home.html", defaultLayout))
	server := Server{&db,mux, oauthCfg,cookieStore, templates}
	server.routes()
	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//https://neoteric.eu/blog/how-to-serve-static-files-with-golang/
	s.router.ServeHTTP(w, r)
}

func NewIDP(config *IdpConfig) (*oauth2.Config, error){

	switch config.Idp  {
		case GitHubProvider:
			return &oauth2.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  githubAuthorizeUrl,
					TokenURL: githubTokenUrl,
				},
				RedirectURL: redirectUrl,
				Scopes:      []string{"repo"},
			},nil
			break

		default:
			log.Fatal().Msgf("not implemented %v",config.Idp)
			return nil, nil
	}
	return nil, nil
}



func LoadConfig() (*IdpConfig, error) {
	var config IdpConfig

	b, err := ioutil.ReadFile(defaultConfigFile)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	config.Idp = GitHubProvider

	return &config, nil
}


