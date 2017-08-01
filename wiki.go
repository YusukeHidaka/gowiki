package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
)

type Page struct {
  Title string
  Body  []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

// ReadFileのエラーを考えない実装
// func loadPage(title string) *Page {
//     filename := title + ".txt"
//     body, _ := ioutil.ReadFile(filename)
//     return &Page{Title: title, Body: body}
// }

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

// エラー処理をしていない例
// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
//     t, _ := template.ParseFiles(tmpl + ".html")
//     t.Execute(w, p)
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
//     t, err := template.ParseFiles(tmpl + ".html")
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     err = t.Execute(w, p)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//     }
// }

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// テンプレートを使わない
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//   title := r.URL.Path[len("/view/"):]
//   p, _ := loadPage(title)
//   fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
// }

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

// テンプレートを使わない
// func editHandler(w http.ResponseWriter, r *http.Request) {
//   title := r.URL.Path[len("/edit/"):]
//   p, err := loadPage(title)
//   if err != nil {
//       p = &Page{Title: title}
//   }
//   fmt.Fprintf(w, "<h1>Editing %s</h1>"+
//     "<form action=\"/save/%s\" method=\"POST\">"+
//     "<textarea name=\"body\">%s</textarea><br>"+
//     "<input type=\"submit\" value=\"Save\">"+
//     "</form>",
//     p.Title, p.Title, p.Body)
// }

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

// エラー処理をしていない
// func saveHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/save/"):]
//     body := r.FormValue("body")
//     //FormValueの返り値がstringなので、byteに変換している
//     p := &Page{Title: title, Body: []byte(body)}
//     p.save()
//     http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }


func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}


// validation
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}
