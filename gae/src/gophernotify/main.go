package gophernotify

import (
	"appengine"
	"fmt"
	"gophernotify/channel"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func init() {
	http.HandleFunc("/", root)
}

// テンプレート
var templates = template.Must(template.ParseGlob("template/*.html"))

// ルートハンドラ
func root(w http.ResponseWriter, r *http.Request) {

	// 新しいコンテキストを作成
	c := appengine.NewContext(r)

	// リクエスト毎にキーを作成する
	client, err := channel.NewClient(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}
	client.Listen(c, "post")

	// チャネル関係の初期化
	channel.Init(client.ClientID)

	// クライアント毎にindexを作る
	urlStr := fmt.Sprintf("/%s", client.ClientID)
	http.HandleFunc(urlStr, index)

	// リダイレクトする
	http.Redirect(w, r, urlStr, http.StatusFound)
}

// インデックス
func index(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// クライアントIDを取得する
	uriArray := strings.Split(r.URL.RequestURI(), "/")
	if uriArray == nil || len(uriArray) <= 0 {
		err := fmt.Errorf("Can not get clientID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf(err.Error())
		return
	}
	clientId, _ := strconv.ParseInt(uriArray[len(uriArray)-1], 10, 64)

	// クライアント情報
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
