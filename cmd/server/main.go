package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/junnotantra/backend-test/internal/handler"
	"github.com/junnotantra/backend-test/internal/repository"
	"github.com/junnotantra/backend-test/internal/service"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/inventory.db"
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal(err)
	}

	repo, err := repository.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	if err := migrate(repo.DB()); err != nil {
		log.Fatal(err)
	}

	itemService := service.NewItemService(repo)
	server := &http.Server{Addr: ":8080", Handler: handler.NewHandler(itemService)}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func migrate(db *sql.DB) error {
	entries, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(entries)
	var current int
	if err := db.QueryRow("PRAGMA user_version").Scan(&current); err != nil {
		return err
	}
	for _, entry := range entries {
		version, err := strconv.Atoi(filepath.Base(entry)[:3])
		if err != nil || version <= current {
			continue
		}
		sqlText, err := fs.ReadFile(migrations, entry)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(sqlText)); err != nil {
			return err
		}
		if _, err := db.Exec("PRAGMA user_version = " + strconv.Itoa(version)); err != nil {
			return err
		}
		current = version
	}
	return nil
}
