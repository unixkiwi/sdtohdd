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
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verbose bool

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

func moveFile(src, dst string) error {
	// check if file exists
	_, err := os.Stat(dst)
	if !errors.Is(err, os.ErrNotExist) {
		slog.Warn("File " + dst + " already exists! Skipping " + src)
		return err
	}

	stat, err := os.Stat(src)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Debug("Opening file " + src)
	inputFile, err := os.Open(src)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer inputFile.Close()

	slog.Debug("Creating file " + dst)
	outputFile, err := os.Create(dst)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer outputFile.Close()

	slog.Debug("Copying perms " + stat.Mode().Perm().String() + " to " + dst)
	err = os.Chmod(dst, stat.Mode().Perm())
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Debug("Copying file " + inputFile.Name() + " to " + outputFile.Name())
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// do not remove
	// err = os.Remove(src)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return err
	// }

	return nil
}

func runSdToHdd(cmd *cobra.Command, args []string) {
	dest := args[len(args)-1]
	sources := args[:len(args)-1]

	dstFileInfo, err := os.Stat(dest)
	if err != nil {
		slog.Error(err.Error())
	}

	if !dstFileInfo.IsDir() {
		slog.Error("Destination must be a directory! (" + dest + ")")
		return
	}

	if !strings.HasSuffix(dest, "/") {
		dest = dest + "/"
	}

	slog.Debug("Sources: " + strings.Join(sources, " "))
	slog.Debug("Destination: " + dest)

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

		if !strings.HasSuffix(source, "/") {
			source = source + "/"
		}

		slog.Debug("Walking Dir " + source)

		filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				slog.Error(err.Error())
				return err
			} else {
				if !entry.IsDir() {
					slog.Debug("Attempting to copy " + path + " to " + dest + entry.Name())
					if err = moveFile(path, dest+entry.Name()); err != nil {
						slog.Error("Failed to copy " + source + entry.Name() + " to " + dest + entry.Name() + "!")
					}
				} else {
					slog.Info("Skipping directory " + path)
				}
				return nil
			}
		})
	}
}
