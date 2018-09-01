package launchpad

import (
	"sync"
	"github.com/kataras/iris"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/context"
)

type Instance struct {
	sync.Mutex
}

func GetAllServers(ctx iris.Context) {
	servers := make([]Server, 0)

	var serverModels []ServerModel
	dbInstance.db.Find(&serverModels)

	for _, server := range serverModels {
		stats := ServerStats{}
		statsModel := &ServerStatsModel{}

		dbInstance.db.Where("server_id = ?", server.ID).Order("recorded_at DESC").First(statsModel)

		if statsModel != nil {
			stats.OnlinePlayers = statsModel.OnlineCount
			stats.RegisteredPlayers = statsModel.RegisteredCount
		}

		servers = append(servers, Server{
			ID:                 server.ServerID,
			DistributionURL:    server.DistributionURL,
			LastUpdated:        server.LastUpdated.Unix(),
			DiscordPresenceKey: server.DiscordPresenceKey,
			IPAddress:          server.IPAddress,
			Name:               server.ServerName,
			Category:           server.ServerCategory,
			Stats:              stats,
		})
	}

	ctx.JSON(servers)
}

func GetServer(ctx iris.Context) {
	serverId := ctx.Params().Get("serverId")

	fmt.Println(serverId)

	var server ServerModel

	if err := dbInstance.db.Where("server_id = ?", serverId).First(&server).Error; err != nil {
		ctx.StatusCode(404)
		ctx.JSON(iris.Map{
			"message": "Server not found",
		})

		return
	}

	ctx.JSON(Server{
		ID:                 server.ServerID,
		DistributionURL:    server.DistributionURL,
		LastUpdated:        server.LastUpdated.Unix(),
		DiscordPresenceKey: server.DiscordPresenceKey,
		IPAddress:          server.IPAddress,
		Name:               server.ServerName,
		Category:           server.ServerCategory,
	})
}

func GetServerMods(ctx iris.Context) {
	serverId := ctx.Params().Get("serverId")

	fmt.Println(serverId)

	var server ServerModel

	if err := dbInstance.db.Where("server_id = ?", serverId).First(&server).Error; err != nil {
		ctx.StatusCode(404)
		ctx.JSON(iris.Map{
			"message": "Server not found",
		})

		return
	}

	mods := make([]ModInfo, 0)

	var modModels []ModModel

	dbInstance.db.Model(&server).
		Where("branch = ?", "STABLE").
		Order("id DESC").
		Related(&modModels, "Mods")

	for _, mod := range modModels {
		files := make([]ModFile, 0)
		requiredMods := make([]string, 0)

		if mod.Files != "" {
			json.NewDecoder(strings.NewReader(mod.Files)).Decode(&files)
		}

		if mod.RequiredMods != "" {
			json.NewDecoder(strings.NewReader(mod.RequiredMods)).Decode(&requiredMods)
		}

		mods = append(mods, ModInfo{
			ID:           mod.ModID,
			CreatedAt:    mod.CreatedAt.Unix(),
			UpdatedAt:    mod.UpdatedAt.Unix(),
			Description:  mod.Description,
			Files:        files,
			RequiredMods: requiredMods,
		})
	}

	ctx.JSON(mods)
}

func CheckForUpdate(ctx iris.Context) {
	if !ctx.URLParamExists("version") {
		ctx.StatusCode(400)
		ctx.Text("Try again...")
		return
	}

	ctx.JSON(CheckForUpdates(ctx.URLParamTrim("version")))
}

func GetChangelog(ctx iris.Context) {
	ctx.Text(GetLatestChangelog())
}

func (i *Instance) StartWebServer() {
	i.Lock()

	app := iris.Default()
	app.Get("/servers", GetAllServers)
	app.Get("/servers/{serverId:string}", GetServer)
	app.Get("/servers/{serverId:string}/mods", GetServerMods)
	app.Get("/launcher/update", CheckForUpdate)
	app.Get("/launcher/changelog", GetChangelog)
	app.Get("/generate_204.php", func(context context.Context) {
		context.StatusCode(204)
		context.Text("You're on the internet!")
	})

	app.Run(iris.Addr(":7888"))

	i.Unlock()
}
