package main

import (
	"github.com/SeijiOmi/points-service/db"
	"github.com/SeijiOmi/points-service/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
