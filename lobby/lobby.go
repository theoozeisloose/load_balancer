package lobby

// Lobby represents a lobby.
type Lobby struct {
	// Name tells the lobby name.
	Name string `json:"name"`
	// MaxPlayers tells the maximum players that can join the lobby.
	MaxPlayers int `json:"maxPlayers"`
	// NumPlayers tells the number of players currently connected to the lobby.
	NumPlayers int `json:"numPlayers"`
	// Host tells the address of the server hosting the lobby.
	Host string `json:"host"`
	// Port tells the port on which the server is hosting the lobby.
	Port int `json:"port"`
}

// NewLobby creates a new lobby with the given parameters.
func NewLobby(name string, maxPlayers int) Lobby {
	return Lobby{Name: name, MaxPlayers: maxPlayers}
}
