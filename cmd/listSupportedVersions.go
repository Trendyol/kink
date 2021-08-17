/*
Copyright © 2021 pe.container <pe.container@trendyol.com>

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
	"errors"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/types"
)

// NewListSupportedVersionsCmd represents the listSupportedVersions command
func NewListSupportedVersionsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-supported-versions",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("you should not provide any arguments")
			}
			tags, err := crane.ListTags(types.ImageRepository)
			if err != nil {
				return fmt.Errorf("reading tags for %s: %v", types.ImageRepository, err)
			}

			for _, tag := range tags {
				fmt.Println(tag)
			}
			return nil
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(NewListSupportedVersionsCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listSupportedVersionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listSupportedVersionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
