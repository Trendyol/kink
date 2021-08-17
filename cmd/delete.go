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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/user"
)

// NewCmdDelete represents the delete command
func NewCmdDelete() *cobra.Command {
	var all, force bool
	var name, namespace string

	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.Client()
			if err != nil {
				return err
			}

			kubeclient := client.CoreV1().Pods(namespace)

			user, err := user.Current()
			if err != nil {
				return err
			}

			hostname, err := os.Hostname()
			if err != nil {
				return err
			}

			options := metav1.DeleteOptions{}
			if force {
				gracePeriodSeconds := int64(0)
				options.GracePeriodSeconds = &gracePeriodSeconds
			}

			ctx := context.TODO()
			if all {
				pods, err := kubeclient.List(ctx, metav1.ListOptions{
					LabelSelector: fmt.Sprintf("runned-by=%s", fmt.Sprintf("%s_%s", user.Username, hostname)),
				})

				if err != nil {
					return err
				}

				for _, pod := range pods.Items {
					fmt.Printf("Deleting %s \n", pod.Name)
					if err := kubeclient.Delete(ctx, pod.Name, options); err != nil {
						return err
					}
				}
				return nil
			}

			if err := kubeclient.Delete(ctx, name, options); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "All pods")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Target namespace")
	cmd.PersistentFlags().StringVarP(&name, "name", "", "", "Target pod name")
	cmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force delete")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewCmdDelete())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
