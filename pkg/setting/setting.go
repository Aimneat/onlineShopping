package setting

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type server struct {
	RunMode      string
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// var ServerSetting = &Server{}

type datasource struct {
	DriverName string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
	Charset    string
	Loc        string
}

// var DatabaseSetting = &Database{}

type app struct {
	ProductPageSize int
	OrderPageSize   int
}

type jwt struct {
	JwtSecret string
}

type AllConfig struct {
	Server     server
	Datasource datasource
	App        app
	Jwt        jwt
}

var TotalConfig AllConfig

func Setup() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s /n", err))
	}
	viper.Unmarshal(&TotalConfig)
}
