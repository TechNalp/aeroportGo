package redis

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Conn redis.Conn

func ConnectRedis(host string) redis.Conn {
	fmt.Print("Tentative de connexion à REDIS: " + host + "...")
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
	return conn
}

func DisconnetRedis(conn redis.Conn) {
	if conn == nil {
		log.Fatal("Impossible de fermer la connexion à REDIS")
	}
	conn.Close()
}
