package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Character struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type ListeCharacters struct {
	Results []Character `json:"results"`
}

func accueil(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("erreur de chargement du template : %v", err)
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	character := getCharacter()

	if err := tmpl.Execute(w, character); err != nil {
		log.Printf("erreur de rendu du template : %v", err)
	}
}

func getCharacter() []Character {
	resp, err := http.Get("https://rickandmortyapi.com/api/character")
	if err != nil {
		fmt.Println("erreur :", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("erreur API :", resp.Status)
		return nil
	}

	var liste ListeCharacters
	if err := json.NewDecoder(resp.Body).Decode(&liste); err != nil {
		fmt.Println("erreur de décodage :", err)
		return nil
	}

	for i := range liste.Results {
		liste.Results[i].ID = i + 1
	}

	return liste.Results
}

func main() {
	http.HandleFunc("/", accueil)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
