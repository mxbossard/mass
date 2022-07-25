/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"

	"mby.fr/mass/internal/workspace"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up <resourceExpr>",
	Short: "up resources",
	Long:  ``,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspace.UpResources(args)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")
	upCmd.PersistentFlags().BoolVarP(&workspace.NoCacheBuild, "no-cache", "", false, "Disable build cache")
	upCmd.PersistentFlags().BoolVarP(&workspace.ForceBuild, "build", "b", false, "Force build")
	upCmd.PersistentFlags().BoolVarP(&workspace.ForcePull, "pull", "p", false, "Force pull")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
