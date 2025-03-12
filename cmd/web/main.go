package main

import (
	"database/sql"
	"crypto/tls"
	"log"
	"net/http"
	"flag"
	"os"
	"html/template"
	"time"
	"github.com/joho/godotenv"
	"github.com/konnenl/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	users *models.UserModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}

func main(){
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := godotenv.Load(); err != nil {
        errorLog.Fatal(err)
    }
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == ""{
		errorLog.Fatal("No DATABASE_URL in .env")
	}
	dsn := flag.String("dsn", dbURL, "Postgresql data sorce name")

	db, err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil{
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true


	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	infoLog.Println("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("pgx", dsn)
	if err != nil{
		return nil, err
	}
	if err = db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}