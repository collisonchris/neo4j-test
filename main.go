package main

import (
	"fmt"
	"net/http"
	"time"
	log "github.com/sirupsen/logrus"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	boltLogs "github.com/johnnadratowski/golang-neo4j-bolt-driver/log"

)

type NeoClient struct {
	Pool bolt.DriverPool
}

func NewNeo4jClient(hostname string, port string, user string, password string) (NeoClient, error) {
	pool, err := bolt.NewDriverPool("bolt://"+user+":"+password+"@"+hostname+":"+port, 10)
	boltLogs.SetLevel("info")
	return NeoClient{pool}, err
}

func getConnectionFromPool(pool bolt.DriverPool) bolt.Conn {
	conn, err := pool.OpenPool()
	if err != nil {
		log.Info("error opening pool")
	}
	return conn
}

var neoClient NeoClient

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "heartbeat")
}

func main() {
	go aliveQuery()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func aliveQuery() {
	neoClient, err := NewNeo4jClient("localhost", "7687", "neo4j", "neo4j1")
	for range time.Tick(time.Second * 2 ) {
        conn := getConnectionFromPool(neoClient.Pool)
		_, err = conn.ExecNeo("return 1", nil)
		if(err == nil) {
			log.Info("healthy ping")
		} else {
			log.WithError(err).Error("error with health ping")
		}
		conn.Close()
    }
	
}