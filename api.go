package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/translate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	client *translate.Client
}

func NewApp() App {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(viper.GetString("apiKey")))
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
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/translate", func(r chi.Router) {
			r.Get("/{targetLang}/{str}", a.generateTest)
		})
	})

	log.Println("Web server is available")
	return http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), r)
}

func (a *App) generateTest(w http.ResponseWriter, r *http.Request) {
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

func sendErr(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	http.Error(w, string(message), code)
}
