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
	"os"
	"os/user"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/Trendyol/kink/pkg/kubernetes"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewCmdList represents the list command
func NewCmdList() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all ephemeral cluster",
		Long: `List all ephemeral cluster
		usage: kink list`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.Client()
			if err != nil {
				return err
			}

			if namespace == "" {
				n, _, err := kubernetes.DefaultClientConfig().Namespace()
				if err != nil {
					return err
				}

				namespace = n
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

			pods, err := kubeclient.List(context.TODO(), metav1.ListOptions{
				LabelSelector: fmt.Sprintf("runned-by=%s", fmt.Sprintf("%s_%s", user.Username, hostname)),
			})
			if err != nil {
				return err
			}

			p := printers.NewTablePrinter(printers.PrintOptions{
				Kind:          schema.ParseGroupKind("Pod"),
				WithKind:      true,
				NoHeaders:     false,
				Wide:          true,
				WithNamespace: true,
				ShowLabels:    true,
			})

			for _, pod := range pods.Items {
				_ = p.PrintObj(pod.DeepCopyObject(), os.Stdout)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Target namespace")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewCmdList())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
