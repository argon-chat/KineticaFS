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
}

var args = &params{}

func main() {
	props := parseFlags()
	for i, prop := range props {
		println(i, prop)
	}
	if args.Server {
		router.Run()
	}
}

func parseFlags() []string {
	server := flag.Bool("server", false, "Run as server")
	s := flag.Bool("s", false, "Alias for --server")
	token := flag.String("token", "", "Authorization token")
	scylla := flag.String("scylla", "localhost:9042", "ScyllaDB host:port")
	migrate := flag.Bool("migrate", false, "Run migrations")

	flag.Parse()

	args.Server = *server || *s
	args.Token = *token
	args.ScillaConnection = *scylla
	args.Migrate = *migrate
	return flag.Args()
}
