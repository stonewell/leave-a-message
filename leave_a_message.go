package main

import (
	"os"
	"fmt"
	"strings"

	"net/http"
	"log"
	"github.com/gorilla/mux"

	"gosrc.io/xmpp"
	//"gosrc.io/xmpp/stanza"
)

func LeaveMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := vars["user"]

	password := os.Getenv(strings.ToUpper(user + "_PASSWORD"))
	domain := os.Getenv("XMPP_DOMAIN")
	xmpp_addr := os.Getenv("XMPP_ADDR")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Invalid Request err: %v", err)
		return
        }

	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: xmpp_addr,
		},
		Jid:          user + "@" + domain,
		Credential:   xmpp.Password(password),
		StreamLogger: os.Stdout,
		Insecure:     true,
		// TLSConfig: tls.Config{InsecureSkipVerify: true},
	}

	fmt.Fprintf(w, "Jid:%v\n", config.Jid)
	fmt.Fprintf(w, "Vars:%v\n", vars)
	fmt.Fprintf(w, "Form:%v\n", r.PostForm)
}

func main() {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/{user}", LeaveMessageHandler).Methods(http.MethodPost, http.MethodPut)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
