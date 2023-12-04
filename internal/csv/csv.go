package csv

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DataToLog(values []string) {
	iata := values[0]
	sensorType := values[1]
	value := values[2]
	timestampNotCut := values[3]
	timestamp := strings.Join(strings.Split(timestampNotCut, "-")[:3], "-")
	path := filepath.Join("..", "..", "logs", strings.ToUpper(iata)+"-"+timestamp+".csv")

	if _, err := os.Stat(filepath.Join("..", "..", "logs")); os.IsNotExist(err) {
		mkdirError := os.Mkdir(filepath.Join("..", "..", "logs"), 0644)
		if mkdirError != nil {
			log.Fatal("Impossible de créer le répertoire logs")
		}
	}

	// Créer un fichier si le fichier de log du jour du capteur n'existe pas
	if _, err := os.Stat(path); os.IsNotExist(err) {
		csvFile, err := os.Create(path)
		defer csvFile.Close()
		if err != nil {
			fmt.Printf("Erreur lors de la création du fichier .csv : %s\n", err)
			return
		}
		csvFile.WriteString("IATA,DATE,TYPE MESURE,VALEUR\n")

	}

	// Ouvre le fichier en mode écriture
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Erreur lors de l'ouverture du fichier %s : %s\n", path, err)
		return
	}
	defer file.Close()

	// Écrit dans le fichier
	_, err = file.WriteString(iata + "," + timestampNotCut + "," + sensorType + "," + value + "\n")
	if err != nil {
		fmt.Printf("Erreur lors de l'écriture dans le fichier %s : %s\n", path, err)
		return
	}
	file.Sync()
}
