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
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/kubernetes"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// NewCmdLoad represents the load command
func NewCmdLoad() *cobra.Command {
	var namespace, clusterName string
	var dockerImages []string

	var cmd = &cobra.Command{
		Use:          "load",
		Short:        "Load Docker images into KinD cluster",
		Long:         `It enables to load Docker images into KinD cluster`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a name as an argument")
			}

			nameArg := args[0]

			client, err := kubernetes.Client()
			if err != nil {
				return err
			}

			podClient := client.CoreV1().Pods(namespace)
			ctx := context.TODO()
			_, err = podClient.Get(ctx, nameArg, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("could not get pod: %v", err)
			}

			// Setup the tar path where the images will be saved
			dir, err := TempDir("", "images-tar")
			if err != nil {
				return errors.New("failed to create tempdir")
			}
			defer os.RemoveAll(dir)
			imagesTarPath := filepath.Join(dir, "images.tar")

			err = save(dockerImages, imagesTarPath)
			if err != nil {
				return err
			}

			containerPath := "/tmp/images.tar"
			err = exec.Command("kubectl", "cp", imagesTarPath, fmt.Sprintf("%s/%s:%s", namespace, nameArg, containerPath)).Run()

			if err != nil {
				return err
			}

			result, err := doExec(nameArg, namespace, "kind-cluster", []string{"docker", "load", "-i", containerPath}, nil)
			if err != nil {
				return err
			}

			log.Println(result)

			for _, n := range dockerImages {
				ref, err := name.ParseReference(n)
				args := []string{"kind", "load", "docker-image", ref.Name(), "--name", clusterName, "-v", "8"}
				result, err := doExec(nameArg, namespace, "kind-cluster", args, nil)

				if err != nil {
					return err
				}

				log.Println(result)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Target namespace")
	cmd.Flags().StringVarP(&clusterName, "cluster-name", "", "", "The name for cluster")
	cmd.Flags().StringArrayVarP(&dockerImages, "docker-image", "", []string{}, "The name for Docker image to be load")

	return cmd
}

// save saves images to dest, as in `docker save`
func save(images []string, dest string) error {
	commandArgs := append([]string{"save", "-o", dest}, images...)
	return exec.Command("docker", commandArgs...).Run()
}

// TempDir is like ioutil.TempDir, but more docker friendly
func TempDir(dir, prefix string) (name string, err error) {
	// create a tempdir as normal
	name, err = ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", err
	}
	// on macOS $TMPDIR is typically /var/..., which is not mountable
	// /private/var/... is the mountable equivalent
	if runtime.GOOS == "darwin" && strings.HasPrefix(name, "/var/") {
		name = filepath.Join("/private", name)
	}
	return name, nil
}

func init() {
	rootCmd.AddCommand(NewCmdLoad())
}
