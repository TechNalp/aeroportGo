package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"archi.org/aeroportGo/internal/connection/redis"
)

var connexionName redis.Conn

func setRoutes() {
	http.HandleFunc("/getDatasBetweenTimeValues", getDatasBetweenTimeValues)
	http.HandleFunc("/getAverageOfDay", getAverageValues)
}

func getDatasBetweenTimeValues(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Requête: " + r.URL.Path + "?" + r.URL.RawQuery)
	query := r.URL.Query()

	// Récupération des paramètres de la requête
	airport, present := query["airport"]
	if !present || len(airport) == 0 {
		http.Error(w, "Invalid airport IATA code, please provide valid parameter [airport - INVALID ]", http.StatusBadRequest)
		return
	}
	typeQueryParam, present := query["type"]
	if !present || len(typeQueryParam) == 0 {
		http.Error(w, "Invalid typeCapteur, please provide valid parameter [typeCapteur - INVALID ]", http.StatusBadRequest)
		return
	}
	dateQueryParam1, present := query["date1"]
	if !present || len(dateQueryParam1) == 0 {
		http.Error(w, "Invalid date or time, please provide valid parameter [date1 - INVALID ]", http.StatusBadRequest)
		return
	}
	dateQueryParam2, present := query["date2"]
	if !present || len(dateQueryParam2) == 0 {
		http.Error(w, "Invalid date or time, please provide valid parameter [date2 - INVALID ]", http.StatusBadRequest)
		return
	}

	// Récuperation des paramètres de la requête
	airportParam := airport[0]
	typeParam := typeQueryParam[0]
	dateParam1 := dateQueryParam1[0]
	dateParam2 := dateQueryParam2[0]

	// Conversion des dates en objet time
	layout := "02-01-2006:15-04-05"
	t1, err := time.Parse(layout, dateParam1)
	if err != nil {
		// Gestion de l'erreur si la conversion en objet time échoue
		http.Error(w, "Invalid date or time format, please provide valid parameter [dateParam1 - INVALID ]", http.StatusBadRequest)
		return
	}
	t2, err := time.Parse(layout, dateParam2)
	if err != nil {
		// Gestion de l'erreur si la conversion en objet time échoue
		http.Error(w, "Invalid date or time format, please provide valid parameter [dateParam2 - INVALID ]", http.StatusBadRequest)
		return
	}

	if t1.After(t2) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Date 2 doit être supérieur à Date 1"))
		return
	}

	// Récupération des données en fonction des deux objets time et de l'aeroport
	data := getData_between_twoTimeValues(t1, t2, airportParam, typeParam, w)
	// Envoi de la donnée au client
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(data))
	return
}

func getAverageValues(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Requête: " + r.URL.Path + "?" + r.URL.RawQuery)
	query := r.URL.Query()
	// Récupération des paramètres de la requête
	airport, present := query["airport"]
	if !present || len(airport) == 0 {
		http.Error(w, "Invalid airport IATA code, please provide a code [airport - INVALID ]", http.StatusBadRequest)
		return
	}
	dateQueryParam1, present := query["date"]
	if !present || len(dateQueryParam1) == 0 {
		http.Error(w, "Invalid date or time, please provide a date [date1 - INVALID ]", http.StatusBadRequest)
		return
	}

	// Récuperation des paramètres de la requête
	airportParam := airport[0]
	dateParam1 := dateQueryParam1[0]

	// Conversion des dates en objet time
	layout := "02-01-2006:15-04-05"
	t1, err := time.Parse(layout, dateParam1)
	if err != nil {
		// Gestion de l'erreur si la conversion en objet time échoue
		http.Error(w, "Invalid date or time, please provide valid parameter [dateParam1 - INVALID ]", http.StatusBadRequest)
		return
	}

	// Récupération des données en fonction des deux objets time et de l'aeroport
	data := getAverage_all_types(t1, airportParam, w)
	// Envoi de la donnée au client en tant que json
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func getAverage(t1 time.Time, airportParam string, typeCapteur string) (float64, int) {
	connexionName := redis.ConnectRedis("localhost:6379")
	defer redis.DisconnetRedis(connexionName)
	somme := 0.0
	number := 0
	data, _ := connexionName.Do("HVALS", "airport:"+airportParam+":"+t1.Format("02-01-2006")+":"+typeCapteur)
	if len(fmt.Sprintf("%s", data)) <= 2 {
		return 0.0, 0
	}
	array := strings.Split(fmt.Sprintf("%s", data)[1:][:len(fmt.Sprintf("%s", data))-2], " ")
	for i := 0; i < len(array); i++ {
		parsed, _ := strconv.ParseFloat(fmt.Sprintf("%s", array[i]), 64)
		somme += parsed
		number += 1
	}

	return somme / float64(number), number
}

func getData_between_twoTimeValues(t1 time.Time, t2 time.Time, airportParam string, typeCapteur string, w http.ResponseWriter) string {

	// Initialisation d'un objet json
	result := "{ \"iata\": \"" + airportParam + "\", \"type\": \"" + typeCapteur + "\", \"data\": {"

	for date := t1; date.Before(t2); date = date.Add(24 * time.Hour) {
		formatedDate := date.Format("02-01-2006")

		if date.Format("02-01-2006") != t1.Format("02-01-2006") && date.Format("02-01-2006") != t2.Format("02-01-2006") { // Si le jour en cours ne fait pas partie de celui de la date de début ou de fin, on peut directement mettre toutes les valeurs
			data, _ := connexionName.Do("HGETALL", "airport:"+airportParam+":"+formatedDate+":"+typeCapteur)
			if len(fmt.Sprintf("%s", data)) <= 2 {
				continue
			}
			split := strings.Split(fmt.Sprintf("%s", data)[1:][:len(fmt.Sprintf("%s", data))-2], " ")
			for i := 0; i < (len(split) - 1); i += 2 {
				result += "\"" + formatedDate + ":" + split[i] + "\" : " + split[i+1] + ", \n"
			}
		} else {
			data, _ := connexionName.Do("HGETALL", "airport:"+airportParam+":"+formatedDate+":"+typeCapteur)
			if len(fmt.Sprintf("%s", data)) <= 2 {
				continue
			}
			keys := strings.Split(fmt.Sprintf("%s", data)[1:][:len(fmt.Sprintf("%s", data))-2], " ")
			for _, key := range keys {
				testDate, _ := time.Parse("02-01-2006-15-04-05", formatedDate+"-"+key)
				if testDate.After(t1.Add(-1*time.Second)) && testDate.Before(t2.Add(1*time.Second)) {
					data, _ := connexionName.Do("HGET", "airport:"+airportParam+":"+formatedDate+":"+typeCapteur, key)
					result += "\"" + formatedDate + ":" + key + "\" : " + fmt.Sprintf("%s", data) + ", \n"
				}
			}
		}

	}

	return result[:len(result)-4] + "}}" //On enlève la dernière virgule

}

func getAverage_all_types(t1 time.Time, airportParam string, w http.ResponseWriter) string {

	// Initialisation de la map des moyennes par type de capteur
	averageMap := make(map[string]string)
	// Initilisation d'une deuxieme date 24h apres t1

	// Creation du tableau des types de capteurs
	types := [3]string{"Heat", "Pressure", "Wind"}
	// Boucle sur les types de capteurs
	for _, typeCapteur := range types {
		avergae, _ := getAverage(t1, airportParam, typeCapteur)
		averageMap[typeCapteur] = strconv.FormatFloat(avergae, 'f', 2, 64)
	}
	// Conversion de la map en json
	averageJson, errJson := json.Marshal(averageMap)

	// Gestion de l'erreur si la conversion en json échoue
	if errJson != nil {
		fmt.Println(errJson)
	}

	// Retour du résultat

	return string(averageJson)

}

func main() {

	args := os.Args

	if len(args) != 3 {
		log.Fatal("Veuillez saisir les arguments de lancement:\n - Port à écouter\n - ip:port de la base de donnée redis")
	}

	port := args[1]
	ipAndPortRedis := args[2]

	// Connexion à la base de données REDIS
	connexionName = redis.ConnectRedis(ipAndPortRedis)
	// Définition des routes
	setRoutes()
	// Lancement du serveur
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// Déconnexion de la base de données REDIS
	redis.DisconnetRedis(connexionName)
}
