package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type RicketMorty struct {
	ID   int
	Name string `json:"name"`
}

type ListeRandM struct {
	Result []RicketMorty `json:"results"`
}

func accueil(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("erreur de chargement du template : %v", err)
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	r&m := getRandM()

	if err := tmpl.Execute(w, r&m); err != nil {
		log.Printf("erreur de rendu du template : %v", err)
	}
}

func getRandM() []RicketMorty {
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

	var liste ListeRandM
	if err := json.NewDecoder(resp.Body).Decode(&liste); err != nil {
		fmt.Println("erreur de décodage :", err)
		return nil
	}

	for i := range liste.Result {
		liste.Result[i].ID = i + 1
	}
	return liste.Result
}

func main() {
	http.HandleFunc("/", accueil)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
