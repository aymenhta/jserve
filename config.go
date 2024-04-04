package main

import "flag"

type config struct {
	dbPath string
	port   string
	client string
}

func parseConfig() *config {
	var cfg config

	flag.StringVar(&cfg.dbPath, "db", "./db.json", "the path to the json databse")
	flag.StringVar(&cfg.port, "port", ":4000", "the port the server will be listening to")
	flag.StringVar(&cfg.client, "client", "*", "the base url of the client, helps with CORS")

	flag.Parse()

	return &cfg
}
