package scylla

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

type ScyllaConnection struct {
	Session     *gocql.Session
	Keyspace    string
	Hosts       []string
	Consistency gocql.Consistency
}

func NewScyllaConnection(connectionString string) (*ScyllaConnection, error) {
	hosts := parseConnectionString(connectionString)
	keyspace := viper.GetString("scylla_keyspace")
	if keyspace == "" {
		keyspace = "kineticafs"
	}

	conn := &ScyllaConnection{
		Hosts:       hosts,
		Keyspace:    keyspace,
		Consistency: gocql.LocalQuorum,
	}

	if err := conn.createKeyspaceIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}

	if err := conn.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to ScyllaDB: %w", err)
	}

	log.Printf("Successfully connected to ScyllaDB cluster: %v, keyspace: %s", hosts, keyspace)
	return conn, nil
}

func parseConnectionString(connectionString string) []string {
	if connectionString == "" {
		return []string{"localhost:9042"}
	}

	hosts := strings.Split(connectionString, ",")
	for i, host := range hosts {
		hosts[i] = strings.TrimSpace(host)
		if !strings.Contains(hosts[i], ":") {
			hosts[i] += ":9042"
		}
	}
	return hosts
}

func (sc *ScyllaConnection) createKeyspaceIfNotExists() error {
	cluster := gocql.NewCluster(sc.Hosts...)
	cluster.Timeout = 30 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Consistency = gocql.LocalQuorum
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	tempSession, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create temporary session: %w", err)
	}
	defer tempSession.Close()
	createKeyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s 
		WITH replication = {
			'class': 'SimpleStrategy', 
			'replication_factor': 1
		}`, sc.Keyspace)
	if err := tempSession.Query(createKeyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace %s: %w", sc.Keyspace, err)
	}
	log.Printf("Keyspace '%s' created or already exists", sc.Keyspace)
	return nil
}
func (sc *ScyllaConnection) Connect() error {
	cluster := gocql.NewCluster(sc.Hosts...)
	cluster.Keyspace = sc.Keyspace
	cluster.Timeout = 30 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Consistency = sc.Consistency

	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	cluster.NumConns = 2
	cluster.MaxPreparedStmts = 1000
	cluster.MaxRoutingKeyInfo = 1000

	cluster.ReconnectInterval = 60 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	sc.Session = session
	return nil
}

func (sc *ScyllaConnection) HealthCheck() error {
	if sc.Session == nil {
		return fmt.Errorf("session is nil")
	}

	var result string
	if err := sc.Session.Query("SELECT now() FROM system.local").Scan(&result); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

func (sc *ScyllaConnection) Close() error {
	if sc.Session != nil {
		sc.Session.Close()
		log.Println("ScyllaDB connection closed")
	}
	return nil
}

func (sc *ScyllaConnection) ExecuteQuery(query string, values ...interface{}) error {
	log.Printf("Executing ScyllaDB query: %s, values: %v", query, values)
	return sc.Session.Query(query, values...).Exec()
}

func (sc *ScyllaConnection) GetSession() *gocql.Session {
	return sc.Session
}
