package launchpad

import (
	"sync"
	"time"
	"net/http"
	"net/url"
	"path"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

type ServerInformationResponse struct {
	ServerMessage      string `json:"messageSrv"`
	HomePageUrl        string `json:"homePageUrl"`
	FacebookUrl        string `json:"facebookUrl"`
	DiscordUrl         string `json:"discordUrl"`
	NumberOfRegistered int    `json:"numberOfRegistered"`
	OnlineCount        int    `json:"onlineNumber"`
}

var mutex sync.Mutex

func StartStatsFetcher() {
	timer := time.Tick(1 * time.Minute)

	var serverModels []ServerModel

	for {
		mutex.Lock()
		dbInstance.db.Find(&serverModels)

		for _, server := range serverModels {
			go func(server ServerModel) {
				u, _ := url.Parse(server.IPAddress)
				u.Path = path.Join(u.Path, "GetServerInformation")

				GetLogger().Debug("Fetching server info for ", server.ServerName, " - ", u)

				resp, err := http.Get(u.String())

				if err != nil {
					GetLogger().Error(err)
					dbInstance.db.Model(&server).Update("ServerStatus", "offline")
				} else {
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)

					if err != nil {
						GetLogger().Error(err)
					} else {
						serverInfo := ServerInformationResponse{}
						json.NewDecoder(bytes.NewReader(body)).Decode(&serverInfo)

						GetLogger().Debug("Got server info for ", server.ServerName)

						dbInstance.db.Model(&server).Update("ServerStatus", "online")

						statsModel := ServerStatsModel{}
						statsModel.ServerID = server.ID
						statsModel.Server = server
						statsModel.OnlineCount = serverInfo.OnlineCount
						statsModel.RegisteredCount = serverInfo.NumberOfRegistered
						statsModel.RecordedAt = time.Now()

						dbInstance.db.Create(&statsModel)
					}
				}
			}(server)
		}

		mutex.Unlock()
		<-timer
	}
}
