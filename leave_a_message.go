package main

import (
	"os"
	"fmt"
	"encoding/json"

	"net/http"
	"log"
	"github.com/gorilla/mux"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

type LamUser struct {
	Name string
	Jid string
	Password string
	ToJid string
}

type LamConfig struct {
	ToJid string

	Users []*LamUser
}

func GetConfig(user string) *LamUser {
	config_file := os.Getenv("LAM_CONFIG")

	if config_file == "" {
		return nil
	}

	f, err := os.Open(config_file)

	if err != nil {
		log.Printf("unable to open config file:%v", config_file)
		return nil;
	}

	var config LamConfig

	dec := json.NewDecoder(f)
	err = dec.Decode(&config)

	f.Close()

	if err != nil {
		log.Printf("invalid content of config file:%v", config_file)
		return nil;
	}

	if config.Users == nil {
		log.Printf("no user defined for config file:%v", config_file)
		return nil;
	}

	for _, u := range config.Users {
		if (u.Name == user) {
			if u.ToJid == "" {
				u.ToJid = config.ToJid
			}

			return u
		}
	}

	log.Printf("user:%v not defined for config file:%v", user, config_file)
	return nil
}

func handleMessage(s xmpp.Sender, p stanza.Packet) {
}

func LeaveMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := vars["user"]

	lam_config := GetConfig(user)

	if lam_config == nil {
		log.Printf("%v user config is not found", user)

		http.Error(w, "Invalid Request", 404)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("parse request form data failed: %+v", err)
		http.Error(w, "Invalid Request", 503)
		return
        }

	config := xmpp.Config{
		Jid:          lam_config.Jid,
		Credential:   xmpp.Password(lam_config.Password),
		StreamLogger: os.Stdout,
		Insecure:     false,
		// TLSConfig: tls.Config{InsecureSkipVerify: true},
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	client, err := xmpp.NewClient(&config, router, errorHandler)

	if err != nil {
		log.Printf("new client failed, %+v", err)
	} else {
		if err = client.Connect(); err != nil {
			log.Printf("XMPP connection failed: %+v", err)
		}

		reply := stanza.Message{Attrs: stanza.Attrs{To: lam_config.ToJid}, Body: BuildBody(r, user)}
		err = client.Send(reply)

		if err != nil {
			log.Printf("client send failed, %+v", err)
		}

		client.Disconnect()
	}


        http.Redirect(w, r, "/dashboard", 302)
}

func errorHandler(err error) {
	log.Printf("xmpp error:%+v", err)
}

func BuildBody(r *http.Request, user string) string {
	body := "Messsage from:" + user + "\n"

	for k, v := range r.Form {
		body += fmt.Sprintf("%v:%v\n", k, v)
	}

	return body
}

func main() {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/{user}", LeaveMessageHandler).Methods(http.MethodPost, http.MethodPut)

	addr := os.Getenv("LAM_ADDR")
	port, has_port := os.LookupEnv("LAM_PORT")

	if !has_port {
		port = "8080"
	}

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(addr + ":" + port, r))
}
