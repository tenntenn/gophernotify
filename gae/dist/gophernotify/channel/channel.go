package channel

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// initialize this package
func Init(clientID string) {
	http.HandleFunc(fmt.Sprintf("/%s/response", clientID), response)
	http.HandleFunc(fmt.Sprintf("/%s/post", clientID), post)
}

// client information
type ClientInfo struct {
	// clientID
	ClientID string
	// token
	Token string
}

// create new client and put on datastore
func NewClient(c appengine.Context) (*ClientInfo, error) {

	// key
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "ClientInfo", nil), &ClientInfo{"", ""})

	// create token and clientID
	clientID := fmt.Sprintf("%d", key.IntID())
	tok, err := channel.Create(c, clientID)
	if err != nil {
		return nil, err
	}

	client := &ClientInfo{clientID, tok}
	c.Infof("Regist %#v", client)
	_, err = datastore.Put(c, key, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// get client info by client id from datastore
func GetClient(c appengine.Context, clientID int64) (*ClientInfo, error) {

	c.Infof("Get Client of %d", clientID)

	clientKey := datastore.NewKey(c, "ClientInfo", "", clientID, nil)
	var client ClientInfo
	err := datastore.Get(c, clientKey, &client)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// start listen a post
func (client *ClientInfo) Listen(c appengine.Context, request string) {
	info := CallBackInfo{request, client.ClientID}
	putCallBack(c, info)
}

// callback information
type CallBackInfo struct {
	// callback request
	Request string
	// clientID
	ClientID string
}

// put callback info on datastore.
func putCallBack(c appengine.Context, info CallBackInfo) error {

	// client info is stored?
	intID, _ := strconv.ParseInt(info.ClientID, 10, 64)
	key := datastore.NewKey(c, "ClientInfo", "", intID, nil)
	var client ClientInfo
	err := datastore.Get(c, key, &client)
	if err != nil {
		c.Errorf("%s", err.Error())
		return err
	}

	// put on datastore
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "CallBackInfo", key), &info)
	if err != nil {
		return err
	}

	return nil
}

// send data client callbacks.
func sendCallBack(c appengine.Context, clientID int64, request string, args interface{}) error {

	// clientKey
	clientKey := datastore.NewKey(c, "ClientInfo", "", clientID, nil)

	// sent data
	handler := fmt.Sprintf("on%s", request)
	data := struct {
		Call string      `json:"call"`
		Args interface{} `json:"args"`
	}{
		handler,
		args,
	}
	q := datastore.NewQuery("CallBackInfo").Ancestor(clientKey).Filter("Request=", request)

	// callbacks
	var callbacks []CallBackInfo
	keys, err := q.GetAll(c, &callbacks)
	if err != nil {
		c.Errorf(err.Error())
		return err
	}

	// send
	for i, callback := range callbacks {
		k := keys[i]
		client := callback.ClientID
		c.Infof("Send to %s", client)
		channel.SendJSON(c, client, data)

		// remove request
		err = datastore.Delete(c, k)
		if err != nil {
			return err
		}
	}

	return nil
}

// response of push by channel api from clients.
// client must do response after receiving data.
func response(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clientId := r.FormValue("clientID")
	c := appengine.NewContext(r)
	c.Infof("Response from %s", clientId)
	info := CallBackInfo{"post", clientId}
	putCallBack(c, info)
}

// post 
func post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	c := appengine.NewContext(r)

	// get clientID from form values or URL
	clientIdStr := r.FormValue("clientID")
	clientId, err := strconv.ParseInt(clientIdStr, 10, 64)
	if err != nil {
		c.Errorf("Cannot get clientID from form values caused by (%s).", err.Error())
		reg, _ := regexp.Compile("^/([0-9]+)/post/?$")
		founds := reg.FindStringSubmatch(r.URL.RequestURI())
		if founds == nil || len(founds) < 2 {
			err := fmt.Errorf("Can not get clientID")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			c.Errorf(err.Error())
			return
		}
		clientId, err = strconv.ParseInt(founds[1], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			c.Errorf("Cannot get clientID from URL caused by (%s).", err.Error())
			return
		}
	}

	msg := struct {
		Body string
	}{
		r.FormValue("message"),
	}
	c.Infof("%s", msg)
	sendCallBack(c, clientId, "post", msg)
}
