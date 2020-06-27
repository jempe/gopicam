package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs"
	"golang.org/x/crypto/bcrypt"

	"github.com/jempe/gopicam/pkg/db"
)

type Server struct {
	Db       *db.DB
	Sessions *scs.SessionManager
}

func setSecureHeaders(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")

	if contentType == "html" {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")

	} else {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

	}
}

func (srv *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	setSecureHeaders(w, "json")

	if r.Method != http.MethodPost {
		returnCode405(w, r)
		return
	}

	// initialize server response
	response := make(map[string]interface{})
	response["access"] = "denied"

	// parse Body Data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	if r.PostForm.Get("username") == string(srv.Db.GetConfigValue("username")) {

		log.Println("Login User Found")

		err = bcrypt.CompareHashAndPassword(srv.Db.GetConfigValue("password"), []byte(r.PostForm.Get("password")))

		if err == nil {
			// prepare successful response
			response["access"] = "granted"

			// Renew the session token...
			err = srv.Sessions.RenewToken(r.Context())
			if err != nil {
				returnCode500(w, r)
				return
			}

			// Save the username in the session
			srv.Sessions.Put(r.Context(), "username", r.PostForm.Get("username"))
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, string(responseJSON))
}

// handler of the Preview image
func (srv *Server) PreviewHandler(w http.ResponseWriter, r *http.Request) {
	setSecureHeaders(w, "json")

	if r.Method != http.MethodGet {
		returnCode405(w, r)
		return
	}

	if srv.Sessions.GetString(r.Context(), "username") != string(srv.Db.GetConfigValue("username")) {
		returnCode401(w, r)
		return
	}
}

func returnCode400(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"status\": \"error\"}"))
}

func returnCode401(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("{\"status\": \"error\"}"))
}

func returnCode403(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("{\"status\": \"error\"}"))
}

func returnCode404(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"status\": \"error\"}"))
}

func returnCode405(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("{\"status\": \"error\"}"))
}

func returnCode500(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"status\": \"error\"}"))
}
