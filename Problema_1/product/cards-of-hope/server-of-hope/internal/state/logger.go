package state

import (
	"os"

	"github.com/charmbracelet/log"
)

// Logger é o logger global utilizado para registrar logs no servidor.
var Logger *log.Logger = log.New(os.Stdout)
