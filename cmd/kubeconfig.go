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
	"github.com/Trendyol/kink/pkg/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(NewCmdKubeConfig())
}

func NewCmdKubeConfig() *cobra.Command {
	var (
		name       string
		namespace  string
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Get running cluster's kubeconfig",
		Long: `Get running cluster's kubeconfig
		usage: kink kubeconfig --name kink-demo --namespace default`,
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

			if outputPath != "" {
				outputPath, err = filepath.Abs(outputPath)
				if err != nil {
					return errors.Wrap(err, "failed to get absolute path of output path")
				}

				pathStat, err := os.Stat(outputPath)
				if err != nil {
					return errors.Wrap(err, "failed to get stat output path")
				}

				if pathStat.IsDir() {
					outputPath = fmt.Sprintf("%s/kink-%s-kubeconfig", outputPath, name)
				}
			}

			ctx := context.TODO()
			serviceClient := client.CoreV1().Services(namespace)

			svc, err := serviceClient.Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			kubeconfig, err := doExec(name, namespace, []string{"kubectl", "config", "view", "--minify", "--flatten"})
			if err != nil {
				return err
			}

			hostIP, err := doExec(name, namespace, []string{"sh", "-c", "echo $CERT_SANS"})
			if err != nil {
				return err
			}

			podIP, err := doExec(name, namespace, []string{"sh", "-c", "echo $API_SERVER_ADDRESS"})
			if err != nil {
				return err
			}

			kubeconfig = strings.ReplaceAll(kubeconfig, podIP, hostIP)

			nodePort := svc.Spec.Ports[0].NodePort
			kubeconfig = strings.ReplaceAll(kubeconfig, "30001", fmt.Sprint(nodePort))

			if outputPath != "" {
				os.WriteFile(outputPath, []byte(kubeconfig), 0o600)
				fmt.Printf("kubeconfig exported for %s to %s\n", name, outputPath)
			} else {
				fmt.Print(kubeconfig)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Cluster name")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Target namespace")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for kubeconfig")

	return cmd
}
