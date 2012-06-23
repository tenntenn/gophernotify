package channel

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"fmt"
	"strconv"
	"net/http"
)

func Init(clientID string) {
	http.HandleFunc(fmt.Sprintf("/%s/response", clientID), response)
	http.HandleFunc(fmt.Sprintf("/%s/post", clientID), post)
}

// 接続しているクライアントの情報です。
type ClientInfo struct {
	// クライアントID
	ClientID string
	// トークンです。
	Token string
}

// 新しいクライアントを作成します。
func NewClient(c appengine.Context) (*ClientInfo, error) {

	// データストアのキー
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "ClientInfo", nil), &ClientInfo{"", ""})

	// トークンの作成
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

// 指定したIDのクライアントを取得します。
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

func (client *ClientInfo) Listen(c appengine.Context, request string) {
	info := CallBackInfo{request, client.ClientID}
	putCallBack(c, info)
}

// クライアントからのコールバックのリクエスト情報です。
type CallBackInfo struct {
	// コールバックのリクエストを表す文字列
	Request string
	// クライアントのID
	ClientID string
}

// コールバック情報を追加します。
func putCallBack(c appengine.Context, info CallBackInfo) error {

	// クライアント情報があるか？
	intID, _ := strconv.ParseInt(info.ClientID, 10, 64)
	key := datastore.NewKey(c, "ClientInfo", "", intID, nil)
	var client ClientInfo
	err := datastore.Get(c, key, &client)
	if err != nil {
		c.Errorf("%s", err.Error())
		return err
	}

	// データストアに登録する
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "CallBackInfo", nil), &info)
	if err != nil {
		return err
	}

	return nil
}

// コールバックにデータを送ります。
func sendCallBack(c appengine.Context, request string, args interface{}) error {

	// 送るデータ
	handler := fmt.Sprintf("on%s", request)
	data := map[string]interface{}{"call": handler, "args": args}
	q := datastore.NewQuery("CallBackInfo").Filter("Request=", request)

	// 送る先
	var callbacks []CallBackInfo
	keys, err := q.GetAll(c, &callbacks)
	if err != nil {
		return err
	}

	// 送信
	for i, callback := range callbacks {
		k := keys[i]
		client := callback.ClientID
		c.Infof("Send to %s", client)
		channel.SendJSON(c, client, data)

		// リクエストを削除
		err = datastore.Delete(c, k)
		if err != nil {
			return err
		}

		// 最後のリクエストならクライアント情報も削除
		// count, err := datastore.NewQuery("CallBackInfo").Filter("ClientID=", callback.ClientID).Count(c)
		// if err != nil {
		// 	return err
		// }
		// if count <= 0 {
		// 	intID, _ := strconv.ParseInt(callback.ClientID, 10, 64)
		// 	clientKey := datastore.NewKey(c, "ClientInfo", "", intID, nil)
		// 	err = datastore.Delete(c, clientKey)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	return nil
}

// Channel APIのプッシュに対するレスポンス
func response(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clientId := r.FormValue("clientID")
	c := appengine.NewContext(r)
	c.Infof("Response from %s", clientId)
	info := CallBackInfo{"post", clientId}	
	putCallBack(c, info)
}

// メッセージの投稿
func post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	c := appengine.NewContext(r)
	msg := struct{
		Body string
	}{
		r.FormValue("message"),
	}
	c.Infof("%s", msg)
	sendCallBack(c, "post", msg)
}