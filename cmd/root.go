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
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	sources []string
	dest    string

	rootCmd = &cobra.Command{
		Use:   "sdtohdd sourceDir1 sourceDirX destinationDir",
		Short: "sdtohdd is a small CLI tool that helps you copy files from your camera's SD card to your HDD sorted by date.",
		Long:  "sdtohdd is a small CLI tool that helps you copy files from your camera's SD card to your HDD sorted by date.",
		Args:  cobra.MinimumNArgs(2),
		Run:   runSdToHdd,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Toggle verbose output")
}

func runSdToHdd(cmd *cobra.Command, args []string) {
	dest = args[len(args)-1]
	sources = args[:len(args)-1]

	slog.Debug("Sources: " + strings.Join(sources, " "))
	slog.Debug("Destination: " + dest)

	dstFileInfo, err := os.Stat(dest)
	if err != nil {
		slog.Error(err.Error())
	}

	if !dstFileInfo.IsDir() {
		slog.Error("Destination must be a directory! (" + dest + ")")
		return
	}

	for _, source := range sources {
		srcFileInfo, err := os.Stat(source)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if !srcFileInfo.IsDir() {
			slog.Error("Source must be a directory! (" + source + ")")
			continue
		}
	}
}
