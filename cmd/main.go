package main

import (
	"denettest/internal/api"
)

func main() {
	api := api.InitAPI()
	api.StartServer()
}
