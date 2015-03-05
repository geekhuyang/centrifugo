package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {

	var port string
	var address string
	var debug bool
	var name string
	var configFile string

	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "Centrifugo",
		Long:  "Centrifuge in GO",
		Run: func(cmd *cobra.Command, args []string) {

			viper.SetDefault("password", "")
			viper.SetDefault("cookie_secret", "cookie_secret")
			viper.SetDefault("api_secret", "api_secret")
			viper.SetDefault("max_channel_length", 255)

			viper.SetConfigFile(configFile)
			viper.ReadInConfig()

			viper.SetEnvPrefix("centrifuge")
			viper.BindEnv("structure", "engine", "insecure", "password", "cookie_secret", "api_secret")

			viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			viper.BindPFlag("address", cmd.Flags().Lookup("address"))
			viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))

			fmt.Printf("%v", viper.AllSettings())

			http.Handle("/", http.FileServer(http.Dir("web/")))
			http.Handle("/connection/", newClientConnectionHandler())

			addr := viper.GetString("address") + ":" + viper.GetString("port")
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		},
	}
	rootCmd.Flags().StringVarP(&port, "port", "p", "8000", "port")
	rootCmd.Flags().StringVarP(&address, "address", "a", "localhost", "address")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.json", "path to config file")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "unique node name")
	rootCmd.Execute()
}
