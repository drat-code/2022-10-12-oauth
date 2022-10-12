package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
)

type oauth2ClientStore struct {
}

func (s *oauth2ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return &models.Client{
		ID:     "drat",
		Domain: "http://localhost:3000",
	}, nil
}

func Auth() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			username := r.URL.Query().Get("username")
			if username == "" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "drat_session",
				Value:    username,
				Expires:  time.Now().Add(time.Minute * 60),
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
			})
		})
		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "drat_session",
				MaxAge:   -1,
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
			})
			http.Redirect(w, r, "/", http.StatusFound)
		})
	}
}

func OAuth2(manager oauth2.Manager) func(r chi.Router) {
	return func(r chi.Router) {
		srv := server.NewDefaultServer(manager)
		srv.SetAllowGetAccessRequest(true)
		srv.SetClientInfoHandler(server.ClientFormHandler)
		srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
			log.Println("Internal Error:", err.Error())
			return
		})
		srv.SetResponseErrorHandler(func(re *errors.Response) {
			log.Println("Response Error:", re.Error.Error())
		})
		srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (string, error) {
			cookie, _ := r.Cookie("drat_session")
			if cookie == nil {
				setRedirectURL(r.URL.String(), w)
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return "", nil
			}
			return cookie.Value, nil
		})
		r.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
			err := srv.HandleAuthorizeRequest(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		})
		r.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
			_ = srv.HandleTokenRequest(w, r)
		})
	}
}

func setRedirectURL(url string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "drat_redirect_url",
		Value:    url,
		Expires:  time.Now().Add(time.Minute * 5),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func main() {
	redisOptions, err := redis.ParseURL("redis://redis:6379")
	if err != nil {
		panic(err)
	}

	manager := manage.NewDefaultManager()
	redisClient := redis.NewClient(redisOptions)
	manager.MapTokenStorage(oredis.NewRedisStoreWithCli(redisClient, "oauth2"))
	manager.MapClientStorage(&oauth2ClientStore{})

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-AirQL-Roles"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: false,
	}))
	r.Route("/auth", Auth())
	r.Route("/oauth2", OAuth2(manager))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
