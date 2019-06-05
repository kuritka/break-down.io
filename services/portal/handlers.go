package portal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *Server) handleAuthCallback(sessionKey string)http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, sessionKey)
		if err != nil {
			fmt.Fprintln(w, "aborted")
			return
		}

		if r.URL.Query().Get("state") != session.Values["state"] {
			fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
			return
		}

		token, err := s.oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
		if err != nil {
			fmt.Fprintln(w, "there was an issue getting your token")
			return
		}

		if !token.Valid() {
			fmt.Fprintln(w, "retreived invalid token")
			return
		}

		client := github.NewClient(s.oauthCfg.Client(oauth2.NoContext, token))

		user, _, err := client.Users.Get(context.Background(), "")
		if err != nil {
			fmt.Println(w, "error getting name")
			return
		}

		session.Values["githubUserName"] = user.Name
		session.Values["githubAccessToken"] = token

		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
	}

}

func (s *Server) handleHome(sessionKey string)http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, sessionKey)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		renderData := map[string]interface{}{}
		if accessToken, ok := session.Values["githubAccessToken"].(*oauth2.Token); ok {
			client := github.NewClient(s.oauthCfg.Client(oauth2.NoContext, accessToken))

			user, _, err := client.Users.Get(context.Background(), "")
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			renderData["github_user"] = user

			var userMap map[string]interface{}
			mapstructure.Decode(user, &userMap)
			renderData["github_user_map"] = userMap
		}

		s.templates["home.html"].ExecuteTemplate(w, "base", renderData)
	}
}

func (s *Server) handleStart(sessionKey string) http.HandlerFunc {
	//x := calledOnce()
	return func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 16)
		rand.Read(b)

		state := base64.URLEncoding.EncodeToString(b)

		session, _ := s.store.Get(r, sessionKey)
		session.Values["state"] = state
		session.Save(r, w)

		url := s.oauthCfg.AuthCodeURL(state)
		http.Redirect(w, r, url, 302)
	}
}


func (s *Server) handleDestroySession(sessionKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, sessionKey)
		if err != nil {
			fmt.Fprintln(w, "aborted")
			return
		}

		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)

	}
}


func (s *Server) handleStream() http.HandlerFunc {
	//todo: implement way back to tell client that server doesn't recognise....
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := s.upgrader.Upgrade(w, r, nil)
		defer socket.Close()
		socket.SetReadLimit(1024)
		//socket.SetReadDeadline(time.Now().Add(60*time.Second))
		//socket.SetPongHandler(func(string) error {	socket.SetReadDeadline(time.Now().Add(60*time.Second));	return nil	})
		if err != nil {
			return
		}
		//socket id opened, now >> go handleRead
		for {
			msgType, msg, err := socket.ReadMessage()
			if err != nil {
				return
			}
			fmt.Printf("%s sent: %s\n", socket.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = socket.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	}
}
