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
	"fmt"
	//"os"
	"strings"
	"log"

	"github.com/spf13/cobra"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/build"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <resourceExpr>",
	Short: "Build resources",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceExpr := strings.Join(args, " ")
		resources, errorz := resources.ResolveExpression(resourceExpr)
		if errorz.GotError() {
			display := display.New()
			display.Display(errorz)
			//os.Exit(1)
		}

		for _, res := range resources {
			builder, err := build.New(res)
			if err != nil {
				log.Fatal(err)
			}
			err = builder.Build()
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Build finished")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
