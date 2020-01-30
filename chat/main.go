package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/Oyekunle-Mark/silent-whisperer/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()

	gomniauth.SetSecurityKey("Thou shalt not pass")
	gomniauth.WithProviders(
		github.New("1be45b00bf7deab19191", "c5a7188b1840cc1b49af1897ff5722c64d02ff85", "http://localhost:8080/auth/callback/github"),
		google.New("868601420359-dsp8abo079nv98sb8gc2ns66736fa7vj.apps.googleusercontent.com", "yZX30kt7tQtq1ebWSssrLwdH", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	})

	// get the room going
	go r.run()

	log.Println("Starting the web server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
