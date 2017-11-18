package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"load_balancer/lobby"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// status represents game process status.
type status struct {
	// PlayersJoined tells whether at least one player joined the game.
	PlayersJoined bool

	// Pid tells the game process pid.
	Pid int
}

const (
	// defaultHost is the default game host.
	defaultHost = "pylon1.usc.edu"

	// startPort is the starting port for the game instance.
	startPort = 9000
)

var (
	// lobbies holds all the currently available lobbies.
	lobbies map[string]*lobby.Lobby = make(map[string]*lobby.Lobby)

	// statuses holds status information about lobbies.
	statuses map[string]*status = make(map[string]*status)

	// lock synchronizes access to the lobbies cache.
	lock sync.RWMutex = sync.RWMutex{}
)

// createKey creates a new key given the host and the port.
func createKey(host string, port int) string {
	return host + "_" + strconv.Itoa(port)
}

// GetLobby writes the list of lobbies.
func GetLobbies(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	names := make([]string, 0, len(lobbies))
	for name := range lobbies {
		names = append(names, name)
	}
	sort.Strings(names)

	lobbyList := make([]lobby.Lobby, 0, len(lobbies))
	for _, name := range names {
		lobby := lobbies[name]
		lobbyList = append(lobbyList, *lobby)
	}
	json.NewEncoder(w).Encode(lobbyList)
}

// identifyNextPort returns the next usable port.
func identifyNextPort() int {
	port := startPort
	for {
		_, exists := lobbies[createKey(defaultHost, port)]
		if !exists {
			return port
		}
		port++
	}
}

// spawnNewGameServer spawns a new game instance on the given port.
func spawnNewGameServer(port int) (int, error) {
	log.Printf("Trying to spawn game on port %d", port)
	portStr := strconv.Itoa(port)
	cmd := exec.Command("/home/gpstudent/game/Linux.x86_64", "-batchmode", "-nographics", "-server", "-port="+portStr, "-logfile", "log_"+portStr+".out")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	log.Printf("Just ran subprocess %d on port %d", cmd.Process.Pid, port)
	return cmd.Process.Pid, nil
}

// CreateLobby creates a new lobby instance.
func CreateLobby(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var l lobby.Lobby
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		log.Println("Failed to create lobby", err)
		return
	}

	port := identifyNextPort()
	pid, err := spawnNewGameServer(port)
	if err != nil {
		log.Println("Failed to create lobby", err)
		return
	}

	newLobby := &lobby.Lobby{
		Name:       l.Name,
		MaxPlayers: 8,
		NumPlayers: 0,
		Host:       defaultHost,
		Port:       port,
	}
	newStatus := &status{
		PlayersJoined: false,
		Pid:           pid,
	}
	lobbies[createKey(newLobby.Host, newLobby.Port)] = newLobby
	statuses[createKey(newLobby.Host, newLobby.Port)] = newStatus
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

	log.Printf("Received update lobby request for port = %d, numPlayers = %d", lobby.Port, lobby.NumPlayers)

	originalLobby := lobbies[createKey(lobby.Host, lobby.Port)]
	originalLobby.NumPlayers = lobby.NumPlayers
	originalStatus := statuses[createKey(lobby.Host, lobby.Port)]
	if originalLobby.NumPlayers > 0 {
		originalStatus.PlayersJoined = true
	}
}

// initLobbies adds the default set of lobbies to the cache.
func initLobbies() {
	lock.Lock()
	defer lock.Unlock()

	lobby1 := &lobby.Lobby{
		Name:       "localhost lobby",
		MaxPlayers: 8,
		NumPlayers: 0,
		Host:       "localhost",
		Port:       8888,
	}
	status1 := &status{
		PlayersJoined: false,
		Pid:           -1,
	}
	lobbies[createKey(lobby1.Host, lobby1.Port)] = lobby1
	statuses[createKey(lobby1.Host, lobby1.Port)] = status1
}

// reaper tries to reap idle process every 10 seconds.
func reaper() {
	for {
		log.Println("Trying to reap game instances...")
		lock.Lock()
		for key, lobby := range lobbies {
			status := statuses[key]
			if lobby.Host == defaultHost && lobby.NumPlayers == 0 && status.PlayersJoined {
				log.Printf("Trying to reap process with id %d", status.Pid)
				err := syscall.Kill(status.Pid, syscall.SIGKILL)
				if err != nil {
					log.Printf("Failed to reap process with id %d", status.Pid)
					continue
				}
				delete(lobbies, key)
				delete(statuses, key)
				log.Printf("Succesfully reaped process with id %d", status.Pid)
			}
		}
		lock.Unlock()
		time.Sleep(10 * time.Second)
	}
}

// ReapLobby marks a lobby for reaping.
func ReapLobby(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	params := mux.Vars(r)
	port, _ := strconv.Atoi(params["port"])
	lobbies[createKey(defaultHost, port)].NumPlayers = 0
	statuses[createKey(defaultHost, port)].PlayersJoined = true
}

// main function bootstraps the load balancer.
func main() {
	initLobbies()
	go reaper()

	router := mux.NewRouter()
	router.HandleFunc("/lobby", GetLobbies).Methods("GET")
	router.HandleFunc("/lobby", CreateLobby).Methods("POST")
	router.HandleFunc("/lobby", UpdateLobby).Methods("PUT")
	router.HandleFunc("/reap/{port}", ReapLobby).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
