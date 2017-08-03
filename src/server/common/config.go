package common

import (
	"flag"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
}
