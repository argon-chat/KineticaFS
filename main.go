package main

import (
	"log"

	_ "github.com/argon-chat/KineticaFS/docs"
	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/argon-chat/KineticaFS/pkg/router"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	repo, err := repositories.NewApplicationRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	if viper.GetBool("migrate") {
		err = repo.Migrate()
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully")
	}
	if viper.GetBool("server") {
		router.Run(viper.GetInt("port"))
	}
}

func initConfig() {
	viper.SetDefault("server", false)
	viper.SetDefault("token", "")
	viper.SetDefault("scylla", "localhost:9042")
	viper.SetDefault("migrate", false)
	viper.SetDefault("port", 3000)
	viper.SetDefault("database", "scylla")

	pflag.BoolP("server", "s", false, "Run as server")
	pflag.String("token", "", "Authorization token")
	pflag.String("scylla", "localhost:9042", "ScyllaDB host:port")
	pflag.BoolP("migrate", "m", false, "Run migrations")
	pflag.Int("port", 3000, "Server port")
	pflag.String("database", "scylla", "Database backend type (scylla, postgres)")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("KINETICAFS")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.kineticafs")
	viper.AddConfigPath("/etc/kineticafs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults and command-line flags")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
