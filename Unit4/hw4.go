package udacity

import (
	"html/template"
	"net/http"
	"regexp"
)

type user struct {
    UserName,
    Password,
    Verify,
    Email,
    ErrUsr   string
}


func HomeWork4_1(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := ""
	pswd := ""
	verify_pswd := ""
	email := ""

	var valid = func() bool {
		// Compile the expression once, usually at init time.
		// Use raw strings to avoid having to quote the backslashes.
		var validName = regexp.MustCompile("^[a-zA-Z0-9_-]{3,20}$")
		var validPswd = regexp.MustCompile("^.{3,20}$")
		var validEmail = regexp.MustCompile(`^[\S]+@[\S]+\.[\S]+$`)

		return validName.Match([]byte(name)) &&
			validPswd.Match([]byte(pswd)) &&
			(pswd == verify_pswd) &&
			(validEmail.Match([]byte(email)) || (email == ""))
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
		usr := user{name, pswd, verify_pswd, email, ""}
		
		if valid() {
//todo:			url := "http://paychks.appspot.com/udacity/hw4_welcome?name=xyz"
			url := "/udacity/hw4_1/welcome"
			http.SetCookie(w, &http.Cookie{Name: "testCookie", Value: "this is the value"})
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
    		println("not valid")
    		usr.ErrUsr = "test"
			t.Execute(w, usr)
		}
	}
}

func HomeWork4_1_welcome(w http.ResponseWriter, r *http.Request) {
	println("welcome")
	t, err := template.ParseFiles("paychks/handlers/udacity/unit2_welcome.html")
	if err != nil {
		println(err.Error())
		return
	}
	str := r.FormValue("name")
	t.Execute(w, str)
}
