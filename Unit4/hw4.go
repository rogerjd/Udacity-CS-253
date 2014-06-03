/*
Copyright (c) 2014, Roger Demagri
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
    * The name of its contributor may not be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package hw4

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	ErrMsgs map[string]string

	User struct {
		Id int64 `datastore: "-"`
		UserName,
		Password,
		Verify,
		Email,
		Hash string
	}

	FormData struct {
		User User

		Errs ErrMsgs
	}
)

func Tst() {
}

func (usr *User) Key(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "User", "", usr.Id, nil)
}

func getUserByName(ctx appengine.Context, un string) *User {
	q := datastore.NewQuery("User").Filter("UserName=", un)
	var u []User
	k, err := q.GetAll(ctx, &u)
	if err != nil {
		return nil
	}

	if len(u) == 0 {
		return nil
	}
	usr := &u[0]
	usr.Id = k[0].IntID()
	return usr
}

func (usr *User) getSalt() string {
	return strings.Split(usr.Hash, "|")[1]
}

func makeHash(name, pswd, salt string) string {

	var makeSalt = func() string {
		var buf bytes.Buffer
		rand.Seed(time.Now().Unix())

		for i := 0; i < 5; i++ {
			buf.WriteString(string(rand.Intn(52) + 33))
		}

		return buf.String()
	}

	if salt == "" {
		salt = makeSalt()
	}

	hash := sha256.New()
	hash.Write([]byte(name + pswd + salt))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr + "|" + salt

}

func (usr *User) Put(ctx appengine.Context) error {
	id, err := datastore.Put(ctx, usr.Key(ctx), usr)
	if err != nil {
		println(err.Error())
	}
	usr.Id = id.IntID()
	return err
}

//ref: from Id in cookie
func GetUserByID(ctx appengine.Context, id string) *User {
	id2, _ := strconv.ParseInt(id, 10, 64)
	usr := &User{Id: id2}
	err := datastore.Get(ctx, usr.Key(ctx), usr)
	if err != nil {
		return nil
	}
	return usr
}

func HomeWork4_1(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := ""
	pswd := ""
	verify_pswd := ""
	email := ""

	ctx := appengine.NewContext(r)

	var valid = func(errMsgs ErrMsgs) bool {
		// Compile the expression once, usually at init time.
		// Use raw strings to avoid having to quote the backslashes.
		var validName = regexp.MustCompile("^[a-zA-Z0-9_-]{3,20}$")
		var validPswd = regexp.MustCompile("^.{3,20}$")
		var validEmail = regexp.MustCompile(`^[\S]+@[\S]+\.[\S]+$`)

		user := getUserByName(ctx, name)
		if user != nil {
			errMsgs["UserName"] = "user exists "
		}
		nameOk := validName.Match([]byte(name))
		if !nameOk {
			errMsgs["UserName"] += "name too small"
		}
		pswdOk := validPswd.Match([]byte(pswd))
		if !pswdOk {
			errMsgs["Pswd"] += "invalid password "
		}
		pswdVerifyOk := pswd == verify_pswd
		if !pswdVerifyOk {
			errMsgs["Pswd"] += "passwords must match"
		}

		return (nameOk) &&
			(pswdOk) &&
			(pswdVerifyOk) &&
			(validEmail.Match([]byte(email)) || (email == "")) &&
			(user == nil)
	}

	t, err := template.ParseFiles("paychks/handlers/udacity/hw4/unit4_signup.html")
	if err != nil {
		println(err.Error())
		return
	}

	if r.Method == "GET" {
		println("get")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		println("post")

		name = r.FormValue("username")
		pswd = r.FormValue("password")
		verify_pswd = r.FormValue("verify")
		email = r.FormValue("email")
		//		usr := user{name, pswd, verify_pswd, email, ""}
		fd := FormData{}
		fd.Errs = make(map[string]string)

		if valid(fd.Errs) {
			usr := User{UserName: name, Password: pswd}
			usr.Hash = makeHash(usr.UserName, usr.Password, "")
			usr.Put(ctx)

			url := "http://paychks.appspot.com/udacity/hw6/"
//			url := "/udacity/hw6/"
			http.SetCookie(w, &http.Cookie{Name: "user_id",
				Value: strconv.FormatInt(usr.Id, 10) + "|" + usr.Hash})
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			println("not valid")
			fd.User.UserName = name
			fd.User.Password = pswd
			t.Execute(w, fd)
		}
	}
}

func HomeWork4_1_welcome(w http.ResponseWriter, r *http.Request) {
	println("welcome")
	t, err := template.ParseFiles("paychks/handlers/udacity/hw4/unit4_welcome.html")
	if err != nil {
		println(err.Error())
		return
	}
	cookie, er := r.Cookie("user_id")
	if er != nil {
		w.Write([]byte(er.Error()))
		return
	}
	hash := cookie.Value
	println(hash)
	flds := strings.Split(hash, "|")
	usr := GetUserByID(appengine.NewContext(r), flds[0])

	if usr == nil {
		//    	t.Execute(w, "not found")
		url := "http://paychks.appspot.com/udacity/hw5_1/blog/signup"
//		url := "/udacity/hw6/signup"
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		t.Execute(w, usr.UserName)
	}
}

func HomeWork4_1_login(w http.ResponseWriter, r *http.Request) {

	var ctx appengine.Context
	var name, pswd string
	var usr *User

	var valid = func() bool {
		usr = getUserByName(ctx, name)
		if usr == nil {
			return false
		}

		salt := usr.getSalt()
		h := makeHash(name, pswd, salt)
		return usr.Hash == h
	}

	println("login")
	t, err := template.ParseFiles("paychks/handlers/udacity/hw4/unit4_login.html")
	if err != nil {
		println(err.Error())
		return
	}

	if r.Method == "GET" {
		println("get")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		println("post")

		ctx = appengine.NewContext(r)

		name = r.FormValue("username")
		pswd = r.FormValue("password")

		if valid() {
			url := "http://paychks.appspot.com/udacity/hw6/"
//			url := "/udacity/hw6/"
			http.SetCookie(w, &http.Cookie{Name: "user_id",
				Value: strconv.FormatInt(usr.Id, 10) + "|" + usr.Hash})
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			println("not valid")
			t.Execute(w, "invalid login")
		}
	}
}

func HomeWork4_1_logout(w http.ResponseWriter, r *http.Request) {

	println("logout")

	if r.Method == "GET" {
		println("get")
		http.SetCookie(w, &http.Cookie{Name: "user_id",
			Value: "", MaxAge: -1})

		//ref: NOTE hw5
		url := "/udacity/hw6/signup"
		http.Redirect(w, r, url, http.StatusFound)
	}
}
