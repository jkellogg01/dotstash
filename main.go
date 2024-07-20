package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/cmd"
)

func main() {
	if os.Getenv("MODE") == "dev" {
		log.SetLevel(log.DebugLevel)
	}
	cmd.Execute()
}
