package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/isotronic/blog-aggregator/internal/config"
	"github.com/isotronic/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg     *config.Config
	db      *database.Queries
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalln(err)
	}
	dbQueries := database.New(db)

	st := state{cfg: &config, db: dbQueries}
	
	cmds := commands{}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)

	clArgs := os.Args
	if len(clArgs) < 2 {
		log.Println("no command given")
		os.Exit(1)
	}

	name := clArgs[1]
	args := clArgs[2:]

	err = cmds.run(&st, command{name, args})
	if err != nil {
		log.Fatalln(err)
	}
}