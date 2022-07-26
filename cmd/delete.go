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
	"errors"
	"fmt"
	"os"
	"os/user"

	corev1 "k8s.io/api/core/v1"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/Trendyol/kink/pkg/kubernetes"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// NewCmdDelete represents the delete command
func NewCmdDelete() *cobra.Command {
	var all, force, silent bool
	var name, namespace string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Ephemeral cluster could be deleted by delete command",
		Long: `You can delete kink cluster by using delete command
		usage:	kink delete`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if silent && name == "" {
				return errors.New("the pod name must be set when using the silent flag")
			}

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

			podClient := client.CoreV1().Pods(namespace)
			serviceClient := client.CoreV1().Services(namespace)

			ctx := context.TODO()

			currentUser, err := user.Current()
			if err != nil {
				return err
			}

			hostname, err := os.Hostname()
			if err != nil {
				return err
			}

			pods, err := podClient.List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("runned-by=%s", fmt.Sprintf("%s_%s", currentUser.Username, hostname)),
			})

			options := metav1.DeleteOptions{}
			if force {
				gracePeriodSeconds := int64(0)
				options.GracePeriodSeconds = &gracePeriodSeconds
				silent = true
			}

			if all {
				if err != nil {
					return err
				}

				for _, p := range pods.Items {
					err := deletePodAndRelatedService(ctx, p, podClient, options, serviceClient, force)
					if err != nil {
						return err
					}
				}
				return nil
			}
			var podNames []string

			for _, pod := range pods.Items {
				podNames = append(podNames, pod.Name)
			}

			var selectedNames []string
			if name != "" {
				selectedNames = append(selectedNames, name)
			} else if !force {
				prompt := &survey.MultiSelect{
					Message: "What pod do you prefer to delete:",
					Options: podNames,
				}
				_ = survey.AskOne(prompt, &selectedNames)
			}

			for _, podName := range selectedNames {
				p, err := podClient.Get(ctx, podName, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("could not get pod: %v", err)
				}

				err = deletePodAndRelatedService(ctx, *p, podClient, options, serviceClient, silent)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "All pods")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Target namespace")
	cmd.Flags().StringVar(&name, "name", "", "Target pod name")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "Don't ask for confirmation, require --name to be set")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force delete")

	return cmd
}

func deletePodAndRelatedService(ctx context.Context, pod corev1.Pod, podClient v1.PodInterface, options metav1.DeleteOptions, serviceClient v1.ServiceInterface, silent bool) error {
	var deleteConfirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Pod %s and Service %s will be deleted... Do you accept?", pod.Name, pod.Name),
	}
	shouldDelete := isContainersReady(pod)

	if silent {
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
		return nil
	}

	err := survey.AskOne(prompt, &deleteConfirm)
	if err != nil {
		return err
	}

	if deleteConfirm {
		var forceDelete bool
		if !shouldDelete {
			p2 := &survey.Confirm{
				Message: "Pod is not ready yet. Do you want to force delete?",
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
