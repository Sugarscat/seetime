package main

import (
	"net/http"
)

func main() {
	err := http.ListenAndServe(":6060", http.FileServer(http.Dir("build")))
	if err != nil {
		print("错误")
		return
	}
}
