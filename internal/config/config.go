package config

import (
	"cmp"
	"flag"
	"os"
	"strconv"
)

type ConfigStruct struct {
	Host  string
	Port  int
	Debug bool
}

const (
	defaultHost = "localhost"
	defaultPort = 8080
)

func ReadConfig() ConfigStruct {

	var cfg ConfigStruct

	flag.StringVar(&cfg.Host, "host", defaultHost, "Host")
	flag.IntVar(&cfg.Port, "port", defaultPort, "Port")
	flag.BoolVar(&cfg.Debug, "debug", false, "Debug")
	flag.Parse()

	cfg.Host = cmp.Or(os.Getenv("SRV_HOST"), defaultHost)

	if tmp := os.Getenv("SRV_PORT"); tmp != "" {
		cfg.Port, _ = strconv.Atoi(tmp)
		//if err != nil {
		//	log.Println(err.Error())
		//	return cfg
		//	cfg.Port = defaultPort
		//}
		//cfg.Port = port
	}

	return cfg

}
