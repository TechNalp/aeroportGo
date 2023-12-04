package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	mqttConn "archi.org/aeroportGo/internal/connection/mosquitto"
	redisConn "archi.org/aeroportGo/internal/connection/redis"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup
var ipAndPortMosquitto string
var ipAndPortRedis string

func messageHandler(client mqtt.Client, message mqtt.Message) {
	if len(strings.Split(message.Topic(), "/")[1]) > 3 {
		fmt.Println("Code IATA invalide")
		return
	}
	values := strings.Split(string(message.Payload()), ":")
	if values[1] != strings.Split(message.Topic(), "/")[1] {
		fmt.Println("Erreur entre le IATA du topic et le IATA de la donnée")
		return
	}
	conn := redisConn.ConnectRedis(ipAndPortRedis)
	defer redisConn.DisconnetRedis(conn)
	log.Println("Ecriture dans la bdd")
	conn.Do("HSET", "airport:"+values[1]+":"+strings.Join((strings.Split(values[4], "-")[:3]), "-")+":"+values[2], strings.Join((strings.Split(values[4], "-")[3:]), "-"), values[3])
	conn.Do("SET", "airport:"+values[1]+":lastGenerated:"+values[2], values[3]) // Pour pour récupérer facilement les dernières valeurs généré
	// id du capteur : iata : sensor type : value : timestamp
	// airport
}

func main() {
	args := os.Args
	if len(args) != 5 {
		log.Fatalln("Veuillez saisir tous les arguments :\n - hote:port du broker mqtt\n - hote:port de la base de donnée REDIS\n - QOS souhaité\n - Id pour le broker mqqt")
	}
	ipAndPortMosquitto = args[1]
	ipAndPortRedis = args[2]
	qos, err := strconv.Atoi(args[3])
	clientId := args[4]

	if err != nil {
		log.Fatal("Impossible de convertir la QOS en nombre")
	}

	if qos != 0 && qos != 1 && qos != 2 {
		log.Fatal("La QOS doit être 0, 1 ou 2")
	}

	client := mqttConn.Connect("tcp://"+ipAndPortMosquitto, clientId)
	client.Subscribe("airports/#", byte(qos), messageHandler)
	wg.Add(1)
	wg.Wait()
}
