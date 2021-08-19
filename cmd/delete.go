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

	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"github.com/AlecAivazis/survey/v2"
)

// NewCmdDelete represents the delete command
func NewCmdDelete() *cobra.Command {
	var all, force bool
	var name, namespace string

	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Ephemeral cluster could be deleted by delete command",
		Long: `You can delete kink cluster by using delete command
		usage:	kink delete`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.Client()
			if err != nil {
				return err
			}

			podClient := client.CoreV1().Pods(namespace)
			serviceClient := client.CoreV1().Services(namespace)

			currentUser, err := user.Current()
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
				pods, err := podClient.List(ctx, metav1.ListOptions{
					LabelSelector: fmt.Sprintf("runned-by=%s", fmt.Sprintf("%s_%s", currentUser.Username, hostname)),
				})

				if err != nil {
					return err
				}

				for _, pod := range pods.Items {
					err := deletePodAndRelatedService(pod.Name, podClient, ctx, options, serviceClient)
					if err != nil {
						return err
					}
				}
				return nil
			}

			err = deletePodAndRelatedService(name, podClient, ctx, options, serviceClient)
			if err != nil {
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


func deletePodAndRelatedService(name string, podClient v1.PodInterface, ctx context.Context, options metav1.DeleteOptions, serviceClient v1.ServiceInterface) error {


	pod, _ := podClient.Get(ctx, name, metav1.GetOptions{})
	podStatusPhase := string(pod.Status.Phase)

	fmt.Println("Pod status is : ", podStatusPhase)

	switch podStatusPhase {
	case "Pending":
		fmt.Println("Pod has already pending")

	case "Running":

	}


	var deleteConfirm bool
	prompt := &survey.Confirm{
		Message:  fmt.Sprintf("Pod %s and Service %s will be deleted... Do you accept ?", name, name),
	}
	survey.AskOne(prompt, &deleteConfirm)

	if deleteConfirm {
		fmt.Printf("Deleting Pod %s \n", name)
		if err := podClient.Delete(ctx, name, options); err != nil {
			return err
		}

		fmt.Printf("Deleting Service %s \n", name)
		if err := serviceClient.Delete(ctx, name, options); err != nil {return err}


	}
	fmt.Println("Delete operation is discarted")

	return nil
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
