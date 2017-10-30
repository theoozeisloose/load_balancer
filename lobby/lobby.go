package lobby

// Lobby represents a lobby.
type Lobby struct {
	// Name tells the lobby name.
	Name string `json:"name"`
	// MaxPlayers tells the maximum players that can join the lobby.
	MaxPlayers int `json:"maxPlayers"`
	// NumPlayers tells the number of players currently connected to the lobby.
	NumPlayers int `json:"numPlayers"`
	// Server tells the server address hosting the lobby.
	Server string `json:"server"`
}

// NewLobby creates a new lobby with the given parameters.
func NewLobby(name string, maxPlayers int) Lobby {
	return Lobby{Name: name, MaxPlayers: maxPlayers}
}
