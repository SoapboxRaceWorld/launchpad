package launchpad

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"os"
	"time"
	"database/sql/driver"
)

type ServerStatus string

const (
	Unknown ServerStatus = "unknown"
	Online  ServerStatus = "online"
	Offline ServerStatus = "offline"
)

func (u *ServerStatus) Scan(value interface{}) error { *u = ServerStatus(value.([]byte)); return nil }
func (u ServerStatus) Value() (driver.Value, error)  { return string(u), nil }

type ServerModel struct {
	gorm.Model
	ServerID           string
	ServerName         string
	ServerCategory     string
	ServerStatus       ServerStatus       `sql:"not null;type:ENUM('unknown','online','offline');DEFAULT:'unknown'"`
	IPAddress          string
	LastUpdated        time.Time
	DistributionURL    string
	DiscordPresenceKey string
	Mods               []ModModel         `gorm:"foreignkey:ServerID;association_foreignkey:id"`
	Stats              []ServerStatsModel `gorm:"foreignkey:ServerID;association_foreignkey:id"`
}

type ModModel struct {
	gorm.Model
	ModID        string
	Description  string
	Files        string
	RequiredMods string
	Branch       string

	ServerID uint
	Server   ServerModel
}

type ServerStatsModel struct {
	gorm.Model

	RecordedAt      time.Time
	OnlineCount     int
	RegisteredCount int

	ServerID uint
	Server   ServerModel
}

func (ServerModel) TableName() string {
	return "servers"
}

func (ModModel) TableName() string {
	return "mods"
}

func (ServerStatsModel) TableName() string {
	return "server_stats"
}

type DbInstance struct {
	db *gorm.DB
}

var dbInstance *DbInstance

func GetDbInstance() *DbInstance {
	once.Do(func() {
		dbInstance = &DbInstance{}
	})

	return dbInstance
}

func (i *DbInstance) Setup() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))

	if err != nil {
		GetLogger().Fatal(fmt.Sprintf("DB error: %s", err))
		return
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&ServerModel{}, &ModModel{}, &ServerStatsModel{})

	i.db = db
}
