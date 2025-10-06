package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/argon-chat/KineticaFS/docs"
	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/argon-chat/KineticaFS/pkg/router"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var asciiArt = `
  _  __ _          _____  _______  _      
 | |/ /(_)        / ____||__   __|| |     
 | ' /  _  _ __  | |        | |   | |     
 |  <  | || '_ \ | |        | |   | |     
 | . \ | || | | || |____    | |   | |____ 
 |_|\_\|_||_| |_| \_____|   |_|   |______|                                                                               
`

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	if viper.GetBool("bootstrap") {
		bootstrapAdminToken()
		return
	}
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

	serverEnabled := viper.GetBool("server")
	if viper.GetBool("server") {
		// When each service starts, increment the WaitGroup counter
		wg.Add(1)
		// Each service must run in a goroutine and accept the WaitGroup
		// to signal the main thread for a successful shutdown
		go router.Run(ctx, wg, viper.GetInt("port"), repo)
	}

	if serverEnabled {
		// Catch OS signals
		<-quit
		log.Println("Shutdown signal received")
		// Send cancel signal to all active services
		cancel()
		// Wait for all services to finish
		wg.Wait()
	} else {
		cancel()
	}
	// Afterwards, close the database connection
	repo.DB.Close()

	log.Println("Application stopped.")
}

func bootstrapAdminToken() {
	port := viper.GetInt("port")
	url := "http://localhost:" + fmt.Sprint(port) + "/v1/st/bootstrap"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to make bootstrap request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Bootstrap failed: %s\n%s", resp.Status, string(body))
	}
	var token struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AccessKey string `json:"access_key"`
		TokenType int8   `json:"token_type"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		log.Fatalf("Failed to decode bootstrap response: %v", err)
	}
	log.Println("Admin service token created:")
	log.Printf("  ID: %s", token.ID)
	log.Printf("  Name: %s", token.Name)
	log.Printf("  AccessKey: %s", token.AccessKey)
	log.Printf("  TokenType: %d", token.TokenType)
	log.Printf("  CreatedAt: %s", token.CreatedAt)
	log.Printf("  UpdatedAt: %s", token.UpdatedAt)
	log.Println("Store the AccessKey securely; it will not be shown again.")
}

func init() {
	fmt.Println(asciiArt)
	fmt.Println("Welcome to KineticaFS!")
	fmt.Println("Initializing configuration...")
	fmt.Println()
	viper.SetDefault("server", false)
	viper.SetDefault("token", "")
	viper.SetDefault("scylla", "localhost:9042")
	viper.SetDefault("migrate", false)
	viper.SetDefault("port", 3000)
	viper.SetDefault("database", "scylla")
	viper.SetDefault("front-end-path", "/var/www")
	viper.SetDefault("region", "./regions.json")

	pflag.BoolP("server", "s", false, "Run as server")
	pflag.String("token", "", "Authorization token")
	pflag.String("scylla", "localhost:9042", "ScyllaDB host:port")
	pflag.BoolP("migrate", "m", false, "Run migrations")
	pflag.Int("port", 3000, "Server port")
	pflag.StringP("database", "d", "scylla", "Database backend type (scylla, postgres)")
	pflag.BoolP("bootstrap", "b", false, "Bootstrap admin service token (makes HTTP request to /v1/st/bootstrap)")
	pflag.StringP("front-end-path", "f", "/var/www", "Path to front-end folder containing index.html (default: /var/www)")
	pflag.StringP("region", "r", "./regions.json", "Path to regions configuration file (default: ./regions.json)")
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
