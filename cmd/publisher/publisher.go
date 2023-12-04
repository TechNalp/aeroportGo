package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"archi.org/aeroportGo/internal/connection/mosquitto"
	"archi.org/aeroportGo/internal/connection/redis"
	"archi.org/aeroportGo/internal/sensors"
)

/**
	.\publisher.exe [ipAndPort] [sensorId] [iata] [sensorType]

	@ipAndPort: localhost:1883
	@sensorId: 1, 2 or 3
	@iata: cf wikipedia
	@sensorType: Wind, Heat, Pressure
	@QOS : 0, 1, 2
	@delay: temps en ms

	topic => airports/iata/sensorType/
	données => sensorId:iata:sensorType:valeur:YYYY-MM-DD-mm-ss
**/

func main() {
	conn := redis.ConnectRedis("localhost:6379")
	defer redis.DisconnetRedis(conn)
	var value float64
	var sensor *sensors.Sensor

	args := os.Args

	if len(os.Args) != 7 {
		log.Fatal("Veuillez saisir tous les arguments au programme:\n- ip:port du broker mqtt\n- ID du capteur\n- IATA de l'aéroport\n- Type de capteur (Wind, Heat ou Pressure)\n- QOS du capteur\n- Delay entre chaque envoie (en ms)")
	}

	ipAndPort := args[1]
	sensorId, errSensorId := strconv.Atoi(args[2])
	iata := args[3]
	sensorType := args[4]
	qos, errQOS := strconv.Atoi(args[5])
	delay, errDelay := strconv.Atoi(args[6])

	client := mosquitto.Connect("tcp://"+ipAndPort, iata+strconv.Itoa(sensorId))

	if errSensorId != nil {
		fmt.Println(errSensorId)
		log.Fatal(" => Impossible de convertir l'ID du capteur vers un entier")
	}

	if errQOS != nil {
		fmt.Println(errQOS)
		log.Fatal(" => Impossible de convertir le QOS vers un entier")
	}

	if qos != 0 && qos != 1 && qos != 2 {
		log.Fatal("La QOS doit être 0, 1 ou 2")
	}

	if errDelay != nil {
		fmt.Println(errDelay)
		log.Fatal(" => Impossible de convertir le delay vers un entier")
	}

	switch sensorType {
	case "Wind":
		sensor = sensors.NewSensor(sensorId, iata, &sensors.WindSensor)
	case "Heat":
		sensor = sensors.NewSensor(sensorId, iata, &sensors.HeatSensor)
	case "Pressure":
		sensor = sensors.NewSensor(sensorId, iata, &sensors.PressureSensor)
	default:
		log.Fatal("Veuillez choisir Wind, Heat ou Pressure")
	}
	//date := time.Date(2023, 01, 01, 00, 00, 00, 00, time.UTC)
	for /*; date.Format("02-01-2006-15-04-05") != time.Date(2023, 01, 10, 15, 17, 00, 00, time.UTC).Format("02-01-2006-15-04-05"); date = date.Add(time.Duration(delay) * time.Millisecond) */ {
		date := time.Now().Format("02-01-2006-15-04-05")
		value = math.Round(sensor.GenerateNextData()*100) / 100
		// if strings.Split(date.Format("02-01-2006-15-04-05"), "-")[3] == "00" && strings.Split(date.Format("02-01-2006-15-04-05"), "-")[4] == "00" && strings.Split(date.Format("02-01-2006-15-04-05"), "-")[5] == "00" {
		// 	fmt.Println(date.Format("02-01-2006-15-04-05") + " " + strconv.FormatFloat(value, 'f', -1, 64))
		// }

		token := client.Publish(
			"airports/"+iata+"/"+sensorType,
			byte(qos),
			false,
			strconv.Itoa(sensorId)+":"+iata+":"+sensorType+":"+strconv.FormatFloat(value, 'f', -1, 64)+":"+date,
		)
		token.Wait()
		token.Error()
		//conn.Do("SET", "airport:"+iata+":"+date.Format("02-01-2006-15-04-05")+":"+sensorType, strconv.FormatFloat(value, 'f', -1, 64))
		//conn.Do("HSET", "airport:"+iata+":"+date.Format("02-01-2006")+":"+sensorType, strings.Join((strings.Split(date.Format("02-01-2006-15-04-05"), "-")[3:]), "-"), strconv.FormatFloat(value, 'f', -1, 64))

		log.Default().Println("Message envoyé : " + strconv.Itoa(sensorId) + ":" + iata + ":" + sensorType + ":" + strconv.FormatFloat(value, 'f', -1, 64) + ":" + date + "\n" +
			"Sur le topic : " + "airports/" + iata + "/" + sensorType)

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

}
