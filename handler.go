package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	
	if len(command.args) == 0 {
		return fmt.Errorf("no username was given")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Sprintf("username was set to %v", cmd.args[0])

	return nil
}
