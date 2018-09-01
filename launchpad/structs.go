package launchpad

type ModFile struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Server struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	IPAddress          string      `json:"ip_address"`
	LastUpdated        int64       `json:"last_updated"`
	DistributionURL    string      `json:"distribution_url"`
	DiscordPresenceKey string      `json:"discord_presence_key"`
	Category           string      `json:"category"` // examples: OFFICIAL, POWER, COMMUNITY, TESTING
	Stats              ServerStats `json:"stats"`
}

type ModInfo struct {
	ID           string    `json:"id"`
	CreatedAt    int64     `json:"created_at"`
	UpdatedAt    int64     `json:"updated_at"`
	Description  string    `json:"description"`
	Files        []ModFile `json:"files"`
	RequiredMods []string  `json:"required_mods"`
}

type ServerStats struct {
	OnlinePlayers     int `json:"online_players"`
	RegisteredPlayers int `json:"registered_players"`
}

type LauncherUpdate struct {
	DownloadURL string `json:"download_url"`
}

type LauncherDownloadCounts struct {
	Total          int `json:"total"`
	CurrentVersion int `json:"current_version"`
}

type LauncherUpdatePayload struct {
	ClientVersion  string                 `json:"client_version"` // version sent in the "check update" request
	LatestVersion  string                 `json:"latest_version"` // latest GitHub release
	UpdateExists   bool                   `json:"update_exists"`
	Update         LauncherUpdate         `json:"update"`
	DownloadCounts LauncherDownloadCounts `json:"download_counts"`
}

type LauncherUpdatePacket struct {
	Code    int                   `json:"code"`
	Payload LauncherUpdatePayload `json:"payload"`
}
