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

// runnable represents a component that manages its own lifecycle
// and cooperates with the application's graceful shutdown mechanism.
//
// The Run method should:
//  1. Start the component’s main work loop.
//  2. Monitor the provided context for cancellation (ctx.Done()).
//  3. Gracefully complete ongoing operations upon cancellation.
//  4. Call wg.Done() exactly once before returning.
//
// The method must return a non-nil error if the component fails to start
// or encounters an unrecoverable runtime error.
//
// Run is expected to block until either the context is canceled
// or the component finishes its work naturally.
type runnable interface {
	Run(ctx context.Context, wg *sync.WaitGroup) error
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		wg.Add(1)
		migrator := repositories.NewMigrator(repo)
		if err := migrator.Run(ctx, wg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		repo.InitializeRepo(ctx, repo)
	}

	serverEnabled := viper.GetBool("server")
	if serverEnabled {
		port := viper.GetInt("port")
		wg.Add(1)
		go router.NewRouter(repo, port).Run(ctx, wg)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")
	cancel()
	wg.Wait()

	if err := repo.Close(); err != nil {
		log.Printf("Error closing repository: %v", err)
	}
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
