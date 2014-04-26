package udacity

import (
	"html/template"
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"	
	"strings"
	"strconv"
//	"html/template"	
)

type blogEntry struct {
    Id  int64 `datastore:"-"`
    Title,
    Subject    string
    Created    time.Time
    
    //ref: user entry only, do not store in ds
    Error       string `datastore:"-"`
}

func (post *blogEntry) FormattedDate() string {
	return post.Created.Format(time.RFC822)
}
    
const hw3_blog_html = `<!DOCTYPE html>

<html>
  <head>
    <title>blog</title>
    <link rel="stylesheet" type="text/css" href="/static/main.css">
  </head>

  <body>
    <h2 class="main-title">CS 253 Blog</h2>
  
    {{range .}}
    <div class="post">
      <div class="post-heading">
        <div class="post-title">
            {{ .Title }}
          <a href="/blog/{{ .Id }}" style="color: gray">{{ .Title }}</a>
        </div>    
        <div class="post-date">
            {{ .FormattedDate }}
        </div>
      </div>

      <div class="post-content">
        <pre>{{ .Subject }}</pre>
      </div>
    </div>
    {{end}}    
    
    
  </body>

</html>`    


/*
    {{range .}}
        <label class="post-title">{{.Title}}</label>  <label class="post-date">{{.FormattedDate}}</label>
        <hr />
        <label class="post-content">{{.Subject}}</label>
        <br />        
    {{end}}    
*/    

/*
    <h2 class="main-title">CS 253 Blog</h2>
  
    <div class="post">
        <div class="post-heading">
            <a href="/blog/{{ post.key().id() }}" class="post-title">
                {{ post.title }}
            </a>

        <div class="post-date">
            {{ post.created.strftime("%b %d, %Y") }}
        </div>
    </div>

    <div class="post-content">
      {{ post.content.replace("\n","<br>") | safe }}
    </div>
*/


func getBlogs(c appengine.Context) ([]blogEntry, error) {
    q := datastore.NewQuery("Blog").Limit(10)
    var res []blogEntry
    _, err := q.GetAll(c, &res)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func HomeWork3_blog(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("foo").Parse(hw3_blog_html)
	if err != nil {
		println(err.Error())
		return
	}
	
    ctx := appengine.NewContext(r)
	id := r.FormValue("id")
	println("********", id)
	if id == "" {
        blogs, err := getBlogs(ctx)
        if err != nil {
            t.Execute(w, err)
            return
        }
      	t.Execute(w, blogs)	
	} else {
        id2, _ := strconv.Atoi(id)
//        id2 = int64(id2)
        println(err)
	    k := datastore.NewKey(ctx, "Blog", "", int64(id2), nil)
        var b blogEntry
	    datastore.Get(ctx, k, &b)
	    sb := []blogEntry{b}
	    err = t.Execute(w, sb)
	    if err != nil {
      	    println(err.Error())
  	    }
	}
	
//    postBlog(ctx, &blogEntry{"zzz", "erkjdsf", time.Now()})

/*
	t := template.New("base2.tmpl")
	t = t.Funcs(template.FuncMap{"eq": reflect.DeepEqual})
	t, err = t.ParseFiles("paychks/templates/base2.tmpl", "paychks/templates/result.tmpl")
*/	
}

/*
//todo: redirect to this, permalink
func HomeWork3_getOneblog(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("foo").Parse(hw3_blog_html)
	if err != nil {
		println(err.Error())
		return
	}
	
//	id := r.FormValue("id")

//    ctx := appengine.NewContext(r)
//    blogs, err := getBlogs(ctx)
    if err != nil {
        t.Execute(w, err)
        return
    }

}
*/

const hw3_blog_newpost_html = `<!DOCTYPE html>

<html>
  <head>
    <title>new post</title>
    <link rel="stylesheet" type="text/css" href="/static/main.css">
  </head>

  <body>
    <h2>CS 253 Blog</h2>    
    
    <form method="post">
    <label> subject
      <br/>
      <input type="text", name="subject"/>
    <label/>
    <br/>
    
    <label> blog
      <br/>
      <textarea name="content"> </textarea>
    <label/>
    <br/>
    
    <div class="error">{{.Error}}</div>
    
    <input type="submit" name="submit"> </button>
    </form>
    
  </body>

</html>`    


func postBlog(c appengine.Context, b *blogEntry) error {
    k := datastore.NewKey(c, "Blog", "", 0, nil)
    var (
        k2 *datastore.Key
        err error
    )
    if k2, err = datastore.Put(c, k, b); err != nil {
        println("")
        return err
    }
    b.Id = k2.IntID()
    return nil
}


func HomeWork3_blog_newPost(w http.ResponseWriter, r *http.Request) {
    t, err := template.New("foo").Parse(hw3_blog_newpost_html)

    if r.Method == "GET" {
	    if err != nil {
		    println(err.Error())
		    return
	    }
	    t.Execute(w, nil)	
	    return
	} else if r.Method == "POST" {
        var b blogEntry

	    subj := strings.TrimSpace(r.FormValue("subject"))
	    cont := strings.TrimSpace(r.FormValue("content"))
	    b.Title = subj
	    b.Subject = cont
	    if subj == "" || cont == "" {
            b.Error = "need both Subject and Content"	    
    	    t.Execute(w, b)	
    	    return
	    }
	    
        ctx := appengine.NewContext(r)
        b.Created = time.Now()
        postBlog(ctx, &b)
        http.Redirect(w, r, "/udacity/hw3/blog?id=" + strconv.Itoa(int(b.Id)), http.StatusMovedPermanently)
//        http.Redirect(w, r, "http://paychks.appspot.com/udacity/hw3/blog?id=" + strconv.Itoa(int(b.Id)), http.StatusMovedPermanently)
	}
}


func getBlog() {
}

