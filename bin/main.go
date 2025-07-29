package main

import (
	"flag"
	"github.com/argon-chat/KineticaFS/pkg/router"
)

type params struct {
	Server           bool
	Token            string
	ScillaConnection string
	Migrate          bool
	Port             int
}

var args = &params{}

func main() {
	props := parseArgs()
	for i, prop := range props {
		println(i, prop)
	}
	if args.Server {
		router.Run(args.Port)
	}
}

func parseArgs() []string {
	server := flag.Bool("server", false, "Run as server")
	s := flag.Bool("s", false, "Alias for --server")
	token := flag.String("token", "", "Authorization token")
	port := flag.Int("port", 3000, "Server port")
	scylla := flag.String("scylla", "localhost:9042", "ScyllaDB host:port")
	migrate := flag.Bool("migrate", false, "Run migrations")

	flag.Parse()

	args.Server = *server || *s
	args.Token = *token
	args.ScillaConnection = *scylla
	args.Migrate = *migrate
	args.Port = *port
	return flag.Args()
}
