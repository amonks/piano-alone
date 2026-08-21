package game

import (
	"monks.co/piano-alone/baseurl"
	"monks.co/piano-alone/data"
)

type State struct {
	Configuration *Configuration
	Phase         Phase
	Players       map[string]*Player
	Rendition     []byte

	DisklavierIsConnected bool
	ConductorIsConnected  bool
}

func (p *Performance) MIDIFilePath(baseURL baseurl.BaseURL) string {
	return baseURL.Rest(data.PathMIDIFile, "id", p.Configuration.PerformanceID, "filename", p.Configuration.Title+".midi")
}

func NewState() *State {
	return &State{
		Phase:   GamePhaseUninitialized,
		Players: map[string]*Player{},
	}
}

func (s *State) CountConnectedPlayers() int {
	n := 0
	for _, p := range s.Players {
		if p.ConnectionState == ConnectionStateConnected {
			n++
		}
	}
	return n
}

func (s *State) CountStartedTutorials() int {
	n := 0
	for _, p := range s.Players {
		if p.HasStartedTutorial {
			n++
		}
	}
	return n
}

func (s *State) CountCompletedTutorials() int {
	n := 0
	for _, p := range s.Players {
		if p.HasCompletedTutorial {
			n++
		}
	}
	return n
}

func (s *State) CountSubmittedTracks() int {
	n := 0
	for _, p := range s.Players {
		if p.HasSubmitted {
			n++
		}
	}
	return n
}

func StateFromBytes(bs []byte) (*State, error) {
	s, err := decode[State]("state", bs)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *State) Bytes() []byte { return encode(s) }
