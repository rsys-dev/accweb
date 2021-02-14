package server

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/assetto-corsa-web/accweb/cfg"
	"github.com/sirupsen/logrus"
)

var (
	serverList []ServerSettings
	m          sync.Mutex
)

func LoadServerList() {
	logrus.Info("Loading server list...")
	dir, err := ioutil.ReadDir(cfg.Get().ConfigPath)

	if err != nil {
		logrus.WithError(err).Fatal("Error opening config directory to initialize server list")
	}

	for _, entry := range dir {
		if entry.IsDir() {
			if err := loadServerSettings(entry.Name()); err != nil {
				logrus.WithError(err).WithField("name", entry.Name()).Fatal("Error loading server settings")
			}
		}
	}

	logrus.WithField("servers", len(serverList)).Info("Server list loaded successfully")
}

func loadServerSettings(name string) error {
	server := &ServerSettings{Id: parseServerId(name)}

	if err := loadConfigFromFile(&server.Configuration, filepath.Join(cfg.Get().ConfigPath, name, configurationJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Settings, filepath.Join(cfg.Get().ConfigPath, name, settingsJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Event, filepath.Join(cfg.Get().ConfigPath, name, eventJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.EventRules, filepath.Join(cfg.Get().ConfigPath, name, eventRulesJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Entrylist, filepath.Join(cfg.Get().ConfigPath, name, entrylistJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.Bop, filepath.Join(cfg.Get().ConfigPath, name, bopJsonName)); err != nil {
		return err
	}

	if err := loadConfigFromFile(&server.AssistRules, filepath.Join(cfg.Get().ConfigPath, name, assistRulesJsonName)); err != nil {
		return err
	}

	addServer(server)
	return nil
}

func parseServerId(name string) int {
	id, err := strconv.Atoi(name)

	if err != nil {
		logrus.WithError(err).Fatal("Error parsing server ID from directory name")
	}

	return id
}

func loadConfigFromFile(config interface{}, path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func GetServerList(withPasswords bool) []ServerSettings {
	if withPasswords {
		return serverList
	}

	list := make([]ServerSettings, 0, len(serverList))

	for _, server := range serverList {
		server.Settings.Password = ""
		server.Settings.AdminPassword = ""
		server.Settings.SpectatorPassword = ""
		list = append(list, server)
	}

	return list
}

func GetRunningServerList() []ServerSettings {
	list := make([]ServerSettings, 0, len(serverList))

	for _, server := range serverList {
		if server.PID > 0 {
			server.Settings.Password = ""
			server.Settings.AdminPassword = ""
			server.Settings.SpectatorPassword = ""
			server.PID = 0
			server.Configuration = ConfigurationJson{}
			list = append(list, server)
		}
	}

	return list
}

func GetServerById(id int, withPasswords bool) *ServerSettings {
	m.Lock()
	defer m.Unlock()

	for _, server := range serverList {
		if server.Id == id {
			if !withPasswords {
				server.Settings.Password = ""
				server.Settings.AdminPassword = ""
				server.Settings.SpectatorPassword = ""
			}

			return &server
		}
	}

	return nil
}

func setServer(server *ServerSettings) {
	m.Lock()
	defer m.Unlock()

	for i, s := range serverList {
		if s.Id == server.Id {
			serverList[i] = *server
			break
		}
	}
}

func addServer(server *ServerSettings) {
	m.Lock()
	defer m.Unlock()
	serverList = append(serverList, *server)
}

func removeServer(id int) {
	m.Lock()
	defer m.Unlock()

	for i, s := range serverList {
		if s.Id == id {
			serverList = append(serverList[:i], serverList[i+1:]...)
			break
		}
	}
}
