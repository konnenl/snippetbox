package main

import (
	"database/sql"
	"log"
	"net/http"
	"flag"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
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

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
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