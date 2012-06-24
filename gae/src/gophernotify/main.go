package gophernotify

import (
	"appengine"
	"fmt"
	"gophernotify/channel"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"
)

// initialize this program
func init() {
	http.HandleFunc("/", root)
}

// html templates
var templates = template.Must(template.ParseGlob("template/*.html"))

// root handler
func root(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// create token and clientID by a request
	client, err := channel.NewClient(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}
	client.Listen(c, "post")

	// initialize channel package
	channel.Init(client.ClientID)

	// create index handler by a client
	urlStr := fmt.Sprintf("/%s", client.ClientID)
	http.HandleFunc(urlStr, index)

	// redirect to specific index by clientID
	http.Redirect(w, r, urlStr, http.StatusFound)
}

// index handler
func index(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// get clientID from URL
	reg, _ := regexp.Compile("^/([0-9]+)/?$")
	founds := reg.FindStringSubmatch(r.URL.RequestURI())
	if founds == nil || len(founds) < 2 {
		err := fmt.Errorf("Can not get clientID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf(err.Error())
		return
	}
	clientId, err := strconv.ParseInt(founds[1], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("Cannot get clientID from URL caused by (%s).", err.Error())
		return
	}

	// get client info from datastore
	client, err := channel.GetClient(c, clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf(err.Error())
		return
	}

	// index.html
	templateData := struct {
		Token    string
		ClientID string
	}{
		client.Token,
		client.ClientID,
	}
	if err = templates.ExecuteTemplate(w, "index", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
	}
}
