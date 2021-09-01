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
	"github.com/Trendyol/kink/pkg/types"
)

// NewListSupportedVersionsCmd represents the listSupportedVersions command
func NewListSupportedVersionsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-supported-versions",
		Short: "List all supported k8s versions",
		Long:  `You can checkout all supported k8s versions with list-supported-versions flag`,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("you should not provide any arguments")
			}
			tags, err := crane.ListTags(types.NodeImageRepository)
			if err != nil {
				return fmt.Errorf("reading tags for %s: %v", types.NodeImageRepository, err)
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
