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

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump <resourceExpr>",
	Short: "bump resources",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		workspace.BumpResources(args)
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpCmd.PersistentFlags().String("foo", "", "A help for foo")
	bumpCmd.Flags().BoolVarP(&workspace.BumpMajor, "major", "", false, "Bump major version")
	bumpCmd.Flags().BoolVarP(&workspace.BumpMinor, "minor", "", false, "Bump minor version")
	bumpCmd.MarkFlagsMutuallyExclusive("major", "minor")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
