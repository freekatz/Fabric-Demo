package main

import (
	"github.com/1uvu/Fabric-Demo/api/server"
)

func main() {
	apiServer := new(server.Server)
	apiServer.Run(":9999")
}
