package mosquitto

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Connect(brokerURI string, clientId string) mqtt.Client {
	fmt.Println("Tentative de connection (" + brokerURI + ", " + clientId + ")...")
	opts := CreateClientOptions(brokerURI, clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
		if err := token.Error(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Connexion réussi !")
	return client
}

func CreateClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)
	opts.SetClientID(clientId)
	return opts
}
