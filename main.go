package main

import (
	"github.com/sharkbait0402/blog-aggregator/internal/config"
	"github.com/sharkbait0402/blog-aggregator/internal/database"
	"fmt"
	"os"
	"database/sql"
)

import _ "github.com/lib/pq"

func main() {

	cfg, err := config.Read()
	if err!=nil {
		fmt.Errorf("read unsuccessful")
	}

	dbUrl:="postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

	db, err := sql.Open("postgres", dbUrl) 
	if err!= nil {
		fmt.Errorf("open failed")
	}

	dbQueries := database.New(db)

	st := state{}
	st.cfg = &cfg
	st.db = dbQueries

	cmds := commands {
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmd:= command {
		name: os.Args[1],
		args: os.Args[2:],
	}

	cmds.run(&st, cmd)

	cfg, err = config.Read()
		if err!=nil {
			fmt.Errorf("read unsuccessful")
		}

	fmt.Printf("%+v\n", cfg)

}
