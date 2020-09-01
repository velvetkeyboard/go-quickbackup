/*
Copyright © 2020 Ramon Moraes <ramonmoraes8080@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/velvetkeyboard/go-backup/backends"
	"gitlab.com/velvetkeyboard/go-backup/config"
	"gitlab.com/velvetkeyboard/go-backup/schema"
	"gitlab.com/velvetkeyboard/go-backup/utils"
	"gitlab.com/velvetkeyboard/go-backup/zipfile"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "creates a .zip file and uploads it a location",
	Long:  `TODO Add a long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		schemaName, _ := cmd.Flags().GetString("schema")
		locationName, _ := cmd.Flags().GetString("location")
		configFilePath, _ := cmd.Flags().GetString("config")

		if configFilePath == "" {
			utils.LoggerError("file path for configuration is not present")
			os.Exit(0)
		}

		// TODO check if the file exists

		config := new(config.Configuration)
		config.Init(configFilePath)

		// falling back to defaults
		if schemaName == "" {
			schemaName = config.GetDefaultSchemaName()
		}

		if locationName == "" {
			locationName = config.GetDefaultLocationName()
		}

		utils.LoggerInfo(fmt.Sprintf(
			"Backing up \"%s\" schema to the location set \"%s\"",
			schemaName,
			locationName,
		))

		schema := new(schema.Schema)
		schema.Init(config, schemaName)

		// TODO Might be interesting creating the backup on the /tmp
		currTimeStr := utils.GetCurrentISOTimeString()
		zipFileTitle := "quickbackups-" + schemaName + "-" + currTimeStr
		zipFileNameFull := zipFileTitle + ".zip"

		zfile := new(zipfile.ZipFile)
		zfile.Init(zipFileNameFull)
		for _, path := range schema.Files {
			utils.LoggerInfo(fmt.Sprintf("[Zipping] %s", path))
			fileContentBytes, _ := ioutil.ReadFile(path)
			zfile.AppendBytes(zipFileTitle+path, fileContentBytes)
		}

		zfile.Save()

		location := config.Locations[locationName]

		// TODO Add AWS S3 Backend
		// TODO Add Google Drive Backend
		switch location.Backend {
		case "filesystem":
			backend := new(local.BackendLocalFilesystem)
			backend.Init(location.Path)
			backend.Upload(zipFileNameFull)
		default:
			utils.LoggerError(fmt.Sprintf(
				"Backend \"%s\" is not implemented yet :'(",
				location.Backend,
			))
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	uploadCmd.Flags().StringP("schema", "s", "", "Schema Name")
	uploadCmd.Flags().StringP("location", "l", "", "Location Name")
}