package portal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/websocket"
	errorextensions "github.com/kuritka/break-down.io/common/utils"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
	"net/http"
	"time"
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
			s.login = user.Login

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


var inactiveDuration = time.Second * 20
var pingDuration = time.Second * 1
var closeConnectionDuration = time.Hour * 8

var noActionLifeBlocker *time.Ticker
var recognisedChannel chan bool

var pingTicker = time.NewTicker(pingDuration)

func (s *Server) serveWebSockets() http.HandlerFunc {
	//todo: implement way back to tell client that server doesn't recognise....
	return func(w http.ResponseWriter, r *http.Request) {

		socket, err := s.upgrader.Upgrade(w, r, nil)
		recognisedChannel = make(chan bool)
		noActionLifeBlocker = time.NewTicker(inactiveDuration)
		defer noActionLifeBlocker.Stop()
		defer socket.Close()


		socket.SetReadLimit(1024)
		// Close connection automatically after 24 hrs
		socket.SetPongHandler(func(string) error {
			//killing connection when no pong. Pinging from websockets writer......
			socket.SetWriteDeadline(time.Now().Add(closeConnectionDuration))
			socket.SetReadDeadline(time.Now().Add(closeConnectionDuration)); return nil
		})

		socket.SetCloseHandler(func(int, string) error {
			pingTicker.Stop()
			return nil
		})


		if err != nil {
			return
		}

		go s.websocksWriter(socket)
		//socket id opened, now >> go handleRead
		for {
			messageType, msg, err := socket.ReadMessage()
			if messageType == -1 {
				socket.Close()
				return
			}
			errorextensions.LogOnError(err, "unable read message from socket")
			recognisedChannel <- true
			fmt.Printf("%s sent: %s %s\n", socket.RemoteAddr(), string(msg), *s.login)
		}

	}
}

func (s *Server) websocksWriter(socket *websocket.Conn ) {

	var lastReceived = time.Now()
	defer socket.Close()
	defer pingTicker.Stop()

	for {
		select {

			case <- pingTicker.C:
				err := socket.WriteMessage(websocket.PingMessage, []byte{})
				errorextensions.LogOnError(err, "unable to ping client")
				break

			case <- noActionLifeBlocker.C:
				if time.Now().Sub(lastReceived)  > inactiveDuration {
					err := socket.WriteMessage(websocket.TextMessage, []byte("off"))
					errorextensions.LogOnError(err, "noActionLifeBlocker: unable to send message to client")
					break
				}
				break

			case <- recognisedChannel:
				if time.Now().Sub(lastReceived)  > inactiveDuration {
					err := socket.WriteMessage(websocket.TextMessage, []byte("on"))
					errorextensions.LogOnError(err, "unable to send message to client")
				}
				lastReceived = time.Now()
			break
		}
	}
}
