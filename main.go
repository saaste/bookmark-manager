package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saaste/bookmark-manager/auth"
	"github.com/saaste/bookmark-manager/bookmarks"
	"github.com/saaste/bookmark-manager/config"
	"github.com/saaste/bookmark-manager/handlers"
	"github.com/saaste/bookmark-manager/migrations"
	"github.com/saaste/bookmark-manager/notifications"
)

func main() {
	db, err := sql.Open("sqlite3", "bookmarks.db")
	if err != nil {
		log.Fatalf("opening database failed: %v", err)
	}
	defer db.Close()

	log.Println("Running database migrations...")
	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatalf("running migrations failed: %v", err)
	}

	log.Println("Loading application configuration...")
	appConf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("loading application config failed: %v", err)
	}

	go initializeBookmarkChecker(appConf, db)

	auth := auth.NewAuthenticator(appConf)
	handler := handlers.NewHandler(db, appConf, auth)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler.HandleIndex)
	r.Get("/feed", handler.HandleFeed)
	r.Get("/tags/{tags}", handler.HandleTags)
	r.Get("/tags/{tags}/feed", handler.HandleFeed)
	r.Get("/login", handler.HandleLogin)
	r.Post("/login", handler.HandleLogin)
	r.Get("/logout", handler.HandleLogout)
	r.Get("/admin/bookmarks", handler.HandlePrivateBookmarks)
	r.Get("/admin/bookmarks/add", handler.HandleBookmarkAdd)
	r.Post("/admin/bookmarks/add", handler.HandleBookmarkAdd)
	r.Get("/admin/bookmarks/{bookmarkID}", handler.HandleBookmarkEdit)
	r.Post("/admin/bookmarks/{bookmarkID}", handler.HandleBookmarkEdit)
	r.Get("/admin/bookmarks/{bookmarkID}/delete", handler.HandleBookmarkDelete)
	r.Post("/admin/bookmarks/{bookmarkID}/delete", handler.HandleBookmarkDelete)
	r.Get("/api/metadata", handler.HandleAPIMetadata)
	r.Get("/api/tags", handler.HandleAPITags)

	handler.ServeFiles(r, "/assets", http.Dir(fmt.Sprintf("templates/%s/assets", appConf.Template)))

	log.Printf("Server address: http://localhost:%d", appConf.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", appConf.Port), r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Starting server failed: %v", err)
	}
}

func initializeBookmarkChecker(appConfig *config.AppConfig, db *sql.DB) {
	if appConfig.CheckInterval < 1 {
		log.Println("Bookmark check is disabled. To enable, set check_interval value in the config to 1 or more")
		return
	}

	repo := bookmarks.NewSqliteRepository(db)
	checker := bookmarks.NewBookmarkChecker(repo)
	notifier := notifications.NewNotifier(appConfig)

	interval := time.Duration(appConfig.CheckInterval) * time.Hour
	log.Printf("Bookmark check is enabled, running every %d hour(s)\n", int(interval.Hours()))

	if appConfig.CheckRunOnStartup {
		log.Printf("Bookmark check is configured to run on app start. Running...")
		checkBookmarks(appConfig, checker, notifier)
	}

	for range time.Tick(interval) {
		checkBookmarks(appConfig, checker, notifier)
	}
}

func checkBookmarks(appConfig *config.AppConfig, checker *bookmarks.BookmarkChecker, notifier *notifications.Notifier) {
	bmErrors, err := checker.CheckBookbarks()
	if err != nil {
		log.Printf("ERROR: checking bookmarks failed: %v\n", err)
	}

	if len(bmErrors) > 0 {
		message := ""

		for _, bmError := range bmErrors {
			message += fmt.Sprintln(bmError.Title)
			message += fmt.Sprintln(bmError.URL)
			message += fmt.Sprintln(bmError.Message)
			message += fmt.Sprintln()
		}

		err := notifier.SendGotifyMessage(fmt.Sprintf("%s: Invalid Bookmarks", appConfig.SiteName), message)
		if err != nil {
			log.Printf("ERROR: Sending gotify notification failed: %v", err)
		}
	}
}
