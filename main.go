package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"load_balancer/lobby"
	"log"
	"net/http"
	"sync"
)

var (
	// lobbies holds all the currently available lobbies.
	lobbies map[string]*lobby.Lobby = make(map[string]*lobby.Lobby)

	// lock synchronizes access to the lobbies cache.
	lock sync.RWMutex = sync.RWMutex{}
)

// createKey creates a new key given the host and the port.
func createKey(host string, port int) string {
	return host + "_" + string(port)
}

// GetLobby writes the list of lobbies.
func GetLobbies(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	lobbyList := make([]lobby.Lobby, 0, len(lobbies))
	for _, lobby := range lobbies {
		lobbyList = append(lobbyList, *lobby)
	}
	json.NewEncoder(w).Encode(lobbyList)
}

// UpdateLobby updates the lobby information.
func UpdateLobby(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var lobby lobby.Lobby
	err := json.NewDecoder(r.Body).Decode(&lobby)
	if err != nil {
		log.Println("Failed to update lobby", err)
		return
	}
	originalLobby := lobbies[createKey(lobby.Host, lobby.Port)]
	originalLobby.NumPlayers = lobby.NumPlayers
}

// initLobbies adds the default set of lobbies to the cache.
func initLobbies() {
	lock.Lock()
	defer lock.Unlock()

	lobby1 := &lobby.Lobby{
		Name:       "localhost lobby",
		MaxPlayers: 4,
		NumPlayers: 0,
		Host:       "localhost",
		Port:       8888,
	}
	lobby2 := &lobby.Lobby{
		Name:       "pylon1 lobby",
		MaxPlayers: 4,
		NumPlayers: 0,
		Host:       "pylon1.usc.edu",
		Port:       8888,
	}
	lobbies[createKey(lobby1.Host, lobby1.Port)] = lobby1
	lobbies[createKey(lobby2.Host, lobby2.Port)] = lobby2
}

// main function bootstraps the load balancer.
func main() {
	initLobbies()

	router := mux.NewRouter()
	router.HandleFunc("/lobby", GetLobbies).Methods("GET")
	router.HandleFunc("/lobby", UpdateLobby).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}
