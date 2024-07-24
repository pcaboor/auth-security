package main

import (
	"api"
	"server"
)

func main() {

	api.DB = api.NewSQL("./databases/accounts.db")
	defer api.DB.Close()
	api.DB.InitTables()

	server.RunServer()
}
