package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/translate"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var decoder = schema.NewDecoder()

type User struct {
	Username string `json:"username"`
	APIKey   string `json:"apikey"`
}

type App struct {
	client *translate.Client
}

func NewApp() App {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(os.Getenv("APP_GOOGLE_CLOUD_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	app := App{
		client,
	}
	return app
}

func (a *App) serve() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/translate", func(r chi.Router) {
			//r.Get("/{targetLang}/{str}", a.translate)
			r.Post("/{targetLang}/{str}", a.translateWithUserToken)
		})
	})

	log.Println("Web server is available")
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), r)
}

func (a *App) translate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := context.Background()
	targetLang := chi.URLParam(r, "targetLang")
	str := chi.URLParam(r, "str")

	lang, err := language.Parse(targetLang)
	if err != nil {
		sendErr(w, 500, err.Error())
		return
	}

	resp, err := a.client.Translate(ctx, []string{str}, lang, nil)
	if err != nil {
		sendErr(w, 500, err.Error())
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(resp[0].Text))
}

func (a *App) translateWithUserToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := context.Background()
	targetLang := chi.URLParam(r, "targetLang")
	str := chi.URLParam(r, "str")

	err := r.ParseForm()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}

	var user User
	err = decoder.Decode(&user, r.PostForm)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}

	log.Printf("User %s accessed api", user.Username)

	userClient, err := translate.NewClient(ctx, option.WithAPIKey(user.APIKey))
	if err != nil {
		sendErr(w, 500, err.Error())
	}

	lang, err := language.Parse(targetLang)
	if err != nil {
		sendErr(w, 500, err.Error())
		return
	}

	resp, err := userClient.Translate(ctx, []string{str}, lang, nil)
	if err != nil {
		sendErr(w, 500, err.Error())
		return
	}

	userClient.Close()

	w.WriteHeader(200)
	w.Write([]byte(resp[0].Text))
}

func sendErr(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	http.Error(w, string(message), code)
}
