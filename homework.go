/*
Copyright (c) 2014, Roger Demagri
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
    * The name of its contributor may not be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package udacity

import (
	"html/template"
	"math"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

func HomeWork1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Hello, Udacity!"))
}

const unit2_html = `<!DOCTYPE html>

<html>
  <head>
    <title>Unit 2 Rot 13</title>
  </head>

  <body>
    <h2>Enter some text to ROT13:</h2>
    <form method="post">
      <textarea name="text"
                style="height: 100px; width: 400px;">{{.}}</textarea>
      <br>
      <input type="submit">
    </form>
  </body>

</html>`

//ref: ROT13
func HomeWork2(w http.ResponseWriter, r *http.Request) {

	var makeResult = func(inp string) string {
		res := ""
		for _, c := range inp {
			if unicode.IsUpper(c) {
				c -= 64
				c += 13
				c = rune(math.Mod(float64(c), 26) + 64)
			} else if unicode.IsLower(c) {
				c -= 96
				c += 13
				c = rune(math.Mod(float64(c), 26) + 96)
			}
			res += string(c)
		}
		return res
	}

	t, err := template.New("foo").Parse(unit2_html)
	if err != nil {
		println(err.Error())
		return
	}

	res := ""
	if r.Method == "GET" {
		println("get")

	} else if r.Method == "POST" {
		println("post")
		str := r.FormValue("text")
		res = makeResult(str)
		println(res)
	}

	t.Execute(w, res)
}

func HomeWork2_1(w http.ResponseWriter, r *http.Request) {

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

	t, err := template.ParseFiles("paychks/handlers/udacity/unit2_signup.html")
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
		if valid() {
			url := "http://paychks.appspot.com/udacity/hw2_1_welcome?name=xyz"
			url = strings.Replace(url, "xyz", name, 1)
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			t.Execute(w, nil)
		}
	}
}

func HomeWork2_1_welcome(w http.ResponseWriter, r *http.Request) {
	println("welcome")
	t, err := template.ParseFiles("paychks/handlers/udacity/unit2_welcome.html")
	if err != nil {
		println(err.Error())
		return
	}
	str := r.FormValue("name")
	t.Execute(w, str)
}
