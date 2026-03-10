package main

import(
	"net/http"
	"flag"
	"database/sql"
	"log/slog"
	"os"
	_ "github.com/lib/pq"
	"waduhek/internal/models"
)
// go run ./cmd/web -addr=":8765"
type Application struct{
	logger 		*slog.Logger
	DB 			*sql.DB
}

func main(){
	db, err := models.Connectdb()
	addr := flag.String("addr", "8080", "port value of server")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &Application{
		logger: logger,
		DB: db,
	}
	logger.Info("starting server on","addr", *addr)
	defer db.Close()
	err = http.ListenAndServe(*addr, app.Routes())
	logger.Error(err.Error())
	os.Exit(1)
}
