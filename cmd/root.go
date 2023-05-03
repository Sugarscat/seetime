package cmd

import (
	"os"
)

func Run() {
	if len(os.Args) == 1 {
		Start()
	} else {
		arg := os.Args[1]
		switch arg {
		case "admin":
			GetPwd()
		default:
			Help()
		}
	}
}

func init() {

}
