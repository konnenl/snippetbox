package main

import (
	"database/sql"
	"log"
	"net/http"
	"flag"
	"os"
	"html/template"
	"github.com/joho/godotenv"
	"github.com/konnenl/snippetbox/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
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

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	
	infoLog.Println("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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