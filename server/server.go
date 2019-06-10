package server

import "os/exec"

type ServerSettings struct {
	Id  int       `json:"id"`
	PID int       `json:"pid"` // 0 = stopped, else running
	Cmd *exec.Cmd `json:"-"`

	// ACC server configuration files
	Configuration ConfigurationJson `json:"basic"`
	Settings      SettingsJson      `json:"settings"`
	Event         EventJson         `json:"event"`
}

type ConfigurationJson struct {
	UdpPort         int `json:"udpPort"`
	TcpPort         int `json:"tcpPort"`
	MaxClients      int `json:"maxClients"`
	ConfigVersion   int `json:"configVersion"`
	RegisterToLobby int `json:"registerToLobby"`
}

type SettingsJson struct {
	ServerName                 string `json:"serverName"`
	Password                   string `json:"password"`
	AdminPassword              string `json:"adminPassword"`
	TrackMedalsRequirement     int    `json:"trackMedalsRequirement"`
	SafetyRatingRequirement    int    `json:"safetyRatingRequirement"`
	ConfigVersion              int    `json:"configVersion"`
	RacecraftRatingRequirement int    `json:"racecraftRatingRequirement"`
	SpectatorSlots             int    `json:"spectatorSlots"`
	SpectatorPassword          string `json:"spectatorPassword"`
	DumpLeaderboards           int    `json:"dumpLeaderboards"`
	IsRaceLocked               int    `json:"isRaceLocked"`
}

type EventJson struct {
	Track                     string            `json:"track"`
	EventType                 string            `json:"eventType"`
	PreRaceWaitingTimeSeconds int               `json:"preRaceWaitingTimeSeconds"`
	SessionOverTimeSeconds    int               `json:"sessionOverTimeSeconds"`
	AmbientTemp               int               `json:"ambientTemp"`
	TrackTemp                 int               `json:"trackTemp"`
	CloudLevel                float64           `json:"cloudLevel"`
	Rain                      float64           `json:"rain"`
	WeatherRandomness         int               `json:"weatherRandomness"`
	ConfigVersion             int               `json:"configVersion"`
	Sessions                  []SessionSettings `json:"sessions"`
	PostQualySeconds          int               `json:"postQualySeconds"`
	PostRaceSeconds           int               `json:"postRaceSeconds"`
}

type SessionSettings struct {
	HourOfDay              int    `json:"hourOfDay"`
	DayOfWeekend           int    `json:"dayOfWeekend"`
	TimeMultiplier         int    `json:"timeMultiplier"`
	SessionType            string `json:"sessionType"`
	SessionDurationMinutes int    `json:"sessionDurationMinutes"`
}

func (server *ServerSettings) start(cmd *exec.Cmd) {
	server.PID = cmd.Process.Pid
	server.Cmd = cmd
}

func (server *ServerSettings) stop() {
	server.PID = 0
	server.Cmd = nil
}
