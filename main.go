package main

import (
	"gses2.app/api/pkg"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/presentation/http"
)

func main() {
	config.LoadEnv()
	http.StartRouter(pkg.InitServices())
}
