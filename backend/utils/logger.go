package utils

import (
	"log"
	"os"
)

var (

	Log *log.Logger

)

func InitLogger() {
	Log = log.New(os.Stdout, "[Propanalytix]", log.LstdFlags|log.Lshortfile)
}