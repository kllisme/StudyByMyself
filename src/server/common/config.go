package common

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func SetupConfig() {

	var (
		conf = flag.String("conf", "./runtime.config", "config file path")
	)

	flag.Parse()

	viper.SetConfigName(*conf)
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Read runtime config fail:", err.Error())
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	status := viper.New()
	status.AddConfigPath("./")
	status.SetConfigName("resource/status")
	err = status.ReadInConfig()
	if err != nil {
		log.Fatalf("Read status config fail:", err.Error())
	}
	status.WatchConfig()
	status.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	StatusConfig = status

}

var (
	StatusConfig *viper.Viper
)
