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
	corev1 "k8s.io/api/core/v1"
	"log"
	"os"
	"os/user"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
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
		SilenceUsage: true,
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

				for _, p := range pods.Items {
					err := deletePodAndRelatedService(&p, podClient, ctx, options, serviceClient)
					if err != nil {
						return err
					}
				}
				return nil
			} else {
				//TODO: fzf? kink list?
				if name == "" {
					log.Fatalln("you must provide a pod name via '--name'")
				}

				p, err := podClient.Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("could not get pod: %v", err)
				}

				err = deletePodAndRelatedService(p, podClient, ctx, options, serviceClient)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "All pods")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Target namespace")
	cmd.Flags().StringVarP(&name, "name", "", "", "Target pod name")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force delete")

	return cmd
}

func deletePodAndRelatedService(pod *corev1.Pod, podClient v1.PodInterface, ctx context.Context, options metav1.DeleteOptions, serviceClient v1.ServiceInterface) error {
	var deleteConfirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Pod %s and Service %s will be deleted... Do you accept?", pod.Name, pod.Name),
	}
	err := survey.AskOne(prompt, &deleteConfirm)
	if err != nil {
		return err
	}

	if deleteConfirm {
		shouldDelete := isContainersReady(pod)
		var forceDelete bool
		if !shouldDelete {
			p2 := &survey.Confirm{
				Message: fmt.Sprintf("Pod is not ready yet. Do you want to force delete?"),
			}
			err := survey.AskOne(p2, &forceDelete)
			if err != nil {
				return err
			}
		}
		if shouldDelete || forceDelete {
			fmt.Printf("Deleting Pod %s\n", pod.Name)
			if err := podClient.Delete(ctx, pod.Name, options); err != nil {
				return fmt.Errorf("deleting pod: %q", err)
			}

			if shouldDelete {
				fmt.Printf("Deleting Service %s\n", pod.Name)
				if err := serviceClient.Delete(ctx, pod.Name, options); err != nil {
					return fmt.Errorf("deleting service: %q", err)
				}
			}
		}
	} else {
		fmt.Println("Delete operation is discarded")
	}

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
