package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq" // Third party package
	"kurstonvassel.net/quotebox/pkg/models/postgresql"
)

func setUpDB(dsn string) (*sql.DB, error) {
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Test our connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Dependencies (things/variables)
// Dependency Injection (passing)
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	quotes   *postgresql.EmployeeModel
	session  *sessions.Session
}

func main() {
	// Create a command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("EMPLOYEE_DB_DSN"), "PostgreSQL DSN (Data Source Name)")
	secret := flag.String("secret", "p7Mhd+qQamgHsS*+8Tg7mNXtcjvu@egz", "Secret key")
	flag.Parse()
	// Create a logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var db, err = setUpDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		quotes: &postgresql.QuoteModel{
			DB: db,
		},
		session: session,
	}

	// Create a custom web server
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	// Start our server
	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
