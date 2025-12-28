package main

import (
	"github.com/sharkbait0402/blog-aggregator/internal/config"
	"fmt"
	"os"
)

func main() {

	cfg, err := config.Read()
	if err!=nil {
		fmt.Errorf("read unsuccessful")
	}

	st := state{}
	st.cfg = &cfg

	cmds := commands {
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmd:=args[1]

	if cmd == "login" {
		fmt.Printf("logging in... ")
	}

	cfg, err = config.Read()
		if err!=nil {
			fmt.Errorf("read unsuccessful")
		}

	fmt.Printf("%+v\n", cfg)

}
