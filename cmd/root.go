/*
Copyright © 2026 unixkiwi

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string

	rootCmd = &cobra.Command{
		Use:   "sdtohdd",
		Short: "sdtohdd is a small CLI tool that helps you copy files from your camera's SD card to your HDD sorted by date.",
		Long:  "sdtohdd is a small CLI tool that helps you copy files from your camera's SD card to your HDD sorted by date.",
		Run: func(cmd *cobra.Command, args []string) {
			// ran when no args -> run with settings in config
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file (default is $HOME/.config/sdtohdd/config.toml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)

	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")

		viper.AddConfigPath("$HOME/.config/sdtohdd/")
	}

	err := viper.ReadInConfig()
	if err == nil {
		panic(fmt.Errorf("There was an error while reading the config file: %w", err))
	}
}
