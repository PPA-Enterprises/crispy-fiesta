package main

import (
	"flag"
	"pkg/common/config"
	"pkg/api"
)

func main() {
	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	//checkErr(api.Start(cfg))
	checkErr(api.StreamAPI(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
