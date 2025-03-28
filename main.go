package main

import (
	"fmt"
	"os"

	"github.com/isotronic/blog-aggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	st := state{config: &cfg}

	cmds := commands{}
	cmds.register("login", loginHandler)

	clArgs := os.Args
	if len(clArgs) < 2 {
		fmt.Println("no command given")
		os.Exit(1)
	}

	name := clArgs[1]
	args := clArgs[2:]

	err = cmds.run(&st, command{name, args})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}