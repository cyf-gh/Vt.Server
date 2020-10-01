package http

import (
	err "../err"
	err_code "../err_code"
	orm "../orm"
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
)

// 发布新文章
type PostModel struct {
	Title string
	Text string
	TagIds[] string
}

func NewPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post PostModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByAtk( r )
	err.CheckErr( e )
	glg.Log( account )
	glg.Log( post )
	e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

// 修改文章
type ModifiedPostModel struct {
	Id int64
	Title string
	Text string
	TagIds[] string
}

func ModifyPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post ModifiedPostModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByAtk( r )
	err.CheckErr( e )
	glg.Log( account )
	glg.Log( post )
	e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

// 更改文章，没有文本内容
// 应对流量节约的情况
type ModifiyPostNoTextModel struct {
	Id int64
	Title string
	TagIds[] string
}

func ModifiyPostNoText( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post ModifiyPostNoTextModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByAtk( r )
	err.CheckErr( e )
	e = orm.ModifyPostNoText( post.Id, post.Title, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

func GetPosts( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	user := r.FormValue("user")
	var posts []orm.Post
	var e error
	var a * orm.Account

	if user == "my" {
		a, e = GetAccountByAtk( r )
	} else {
		// 则获取相应名字的所有文章
		a, e = orm.GetAccountByName( user )
	}
	err.CheckErr( e )
	posts, e = orm.GetPostsByOwner( a.Id )
	err.CheckErr( e )
	postsB, e := json.Marshal( posts )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, string(postsB), err_code.MakeHER200 )
}
