/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"

	"github.com/spf13/cobra"

	"mby.fr/mass/internal/settings"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Configure MASS",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		settingsService, err := settings.GetSettingsService()
		if err != nil {
			log.Fatal(err)
		}
		s := settingsService.Settings()
		fmt.Println(s)
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)

	// Here you will define your flags and settingsuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// settingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// settingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

