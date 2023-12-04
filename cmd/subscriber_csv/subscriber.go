package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	mqttConn "archi.org/aeroportGo/internal/connection/mosquitto"
	csvLog "archi.org/aeroportGo/internal/csv"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup
var ipAndPortMosquitto string

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

	// id du capteur : iata : sensor type : value : timestamp
	csvLog.DataToLog(values[1:])

}

func main() {
	args := os.Args
	if len(args) != 4 {
		log.Fatalln("Veuillez saisir tous les arguments :\n - hote:port du broker mqtt\n - QOS souhaité\n - Id pour le broker mqqt")
	}
	ipAndPortMosquitto = args[1]
	qos, err := strconv.Atoi(args[2])
	clientId := args[3]

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
