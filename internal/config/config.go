package config

import (
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() {
	home = os.UserHomeDir()

}
