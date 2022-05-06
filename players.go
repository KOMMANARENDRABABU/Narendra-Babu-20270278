package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Creating a structure of Player .
type Player struct {
	//gorm.Model
	PlayerID int    `json:"playerId" gorm:"primary_key"`
	Name     string `json:"name"`
	Team     string `json:"team"`
	//PlayerScores   []PlayerScores `json:"playerscores" gorm:"foreign_key:ID"`
}

//Creating a structure of PlayerScores.
type PlayerScores struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Match    string `json:"match"`
	Runs     int    `json:"runs"`
	Wickets  int    `json:"wickets"`
	PlayerID int    `json:"playerId" gorm:"foreign_key:"`
}

type PlayerScore struct {
	ID     int            `json:"id"`
	Scores []PlayerScores `json:"scores"`
}

type PlayerScores1 struct {
	PlayerScores []PlayerScore `json:"playerscores"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:Maqb@611361@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("CREATE DATABASE players_db")
	db.Exec("USE players_db")

	// Migration to create tables for Player and Score schema
	db.AutoMigrate(&Player{}, &PlayerScores{})
}

func main() {
	router := mux.NewRouter()
	// Create
	router.HandleFunc("/player", createPlayer).Methods("POST")
	// Create
	router.HandleFunc("/player/{playerId}/score", createPlayerScore).Methods("POST")
	// Read
	router.HandleFunc("/players/{playerId}", getPlayer).Methods("GET")
	// Read-all
	router.HandleFunc("/players", getPlayers).Methods("GET")
	// Read-all
	router.HandleFunc("/playerscores", getPlayerScores1).Methods("GET")
	// Initialize db connection
	initDB()

	log.Fatal(http.ListenAndServe(":8051", router))
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)
	// Creates new player by inserting records in the `players` and `scores` table
	db.Create(&player)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func createPlayerScore(w http.ResponseWriter, r *http.Request) {
	var playerscore PlayerScores
	json.NewDecoder(r.Body).Decode(&playerscore)
	// Creates new player by inserting records in the `players` and `scores` table
	db.Create(&playerscore)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerscore)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var players []Player
	db.Find(&players)
	fmt.Println("abc1")
	json.NewEncoder(w).Encode(players)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputPlayerID := params["playerId"]

	var player Player
	db.First(&player, inputPlayerID)
	json.NewEncoder(w).Encode(player)
}

func getPlayerScores1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var playerscores1 []PlayerScores
	var players []Player
	db.Find(&players)
	db.Find(&playerscores1)
	var Ps1 PlayerScores1
	for _, k := range players {
		var Ps PlayerScore
		Ps.ID = k.PlayerID
		for _, j := range playerscores1 {
			var playerscor PlayerScores
			if k.PlayerID == j.ID {
				playerscor.Match = j.Match
				playerscor.Runs = j.Runs
				playerscor.Wickets = j.Wickets
				Ps.Scores = append(Ps.Scores, playerscor)
			}
		}
		//fmt.Println(Ps)
		Ps1.PlayerScores = append(Ps1.PlayerScores, Ps)
	}
	json.NewEncoder(w).Encode(Ps1)
}
