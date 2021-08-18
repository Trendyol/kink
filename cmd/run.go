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
	"bytes"
	"context"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/kubernetes"
	"gitlab.trendyol.com/platform/base/poc/kink/pkg/types"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/url"
	"os"
	"os/user"
	"strings"
	"time"
)

// NewCmdRun represents the run command
func NewCmdRun() *cobra.Command {
	var k8sVersion, namespace, outputPath string
	var timeout int

	var cmd = &cobra.Command{
		Use:   "run",
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

			uuid := uuid.NewUUID()

			podName := "kind-cluster-" + string(uuid)

			user, err := user.Current()
			if err != nil {
				return err
			}

			hostname, err := os.Hostname()
			if err != nil {
				return err
			}

			kindCluster := &corev1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        podName,
					Annotations: kubernetes.ManagedAnnotations(),
					Labels: map[string]string{
						"runned-by": fmt.Sprintf("%s_%s", user.Username, hostname),
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "varlibdocker",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "libmodules",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/lib/modules",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "kind-cluster",
							Image: types.ImageRepository + ":" + k8sVersion,
							Args: []string{
								"/bin/bash",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "api-server-port",
									HostPort:      0,
									ContainerPort: 30001,
									Protocol:      corev1.Protocol("TCP"),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "API_SERVER_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name: "CERT_SANS",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "varlibdocker",
									MountPath: "/var/lib/docker",
								},
								{
									Name:      "libmodules",
									ReadOnly:  true,
									MountPath: "/lib/modules",
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{
											Type:   intstr.Type(1),
											IntVal: 0,
											StrVal: "api-server-port",
										},
										Scheme: corev1.URIScheme("HTTPS"),
									},
								},
								InitialDelaySeconds: 120,
								TimeoutSeconds:      1,
								PeriodSeconds:       20,
								SuccessThreshold:    1,
								FailureThreshold:    15,
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							SecurityContext: &corev1.SecurityContext{
								Privileged: ptrbool(true),
							},
							Stdin: true,
							TTY:   true,
						},
					},
				},
			}

			// Manage resource
			ctx := context.TODO()
			_, err = kubeclient.Create(ctx, kindCluster, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}

			bar := progressbar.NewOptions(timeout,
				progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(true),
				progressbar.OptionSetWidth(15),
				progressbar.OptionSetDescription(fmt.Sprintf("[cyan][1/3][reset] Creating Pod %s...", podName)),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}))
			err = wait.PollImmediate(time.Second, time.Duration(timeout)*time.Second, func() (done bool, err error) {
				bar.Add(1)
				pod, err := kubeclient.Get(ctx, podName, metav1.GetOptions{})
				if err != nil {
					return false, err
				}

				switch pod.Status.Phase {
				case corev1.PodFailed, corev1.PodSucceeded:
					return false, wait.ErrWaitTimeout
				}

				for _, cs := range pod.Status.ContainerStatuses {
					if cs.Ready && *cs.Started {
						return true, nil
					}
				}

				return false, nil
			})

			if err != nil {
				log.Fatalf("the pod never entered running phase: %v\n", err)
			}

			const tty = false
			req := client.CoreV1().RESTClient().Post().
				Resource("pods").
				Name(podName).
				Namespace(namespace).
				SubResource("exec").
				Param("container", "kind-cluster")

			req.VersionedParams(&corev1.PodExecOptions{
				Container: "kind-cluster",
				Command:   []string{"kubectl", "config", "view", "--minify", "--flatten"},
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       tty,
			}, scheme.ParameterCodec)

			var stdout, stderr bytes.Buffer
			config, err := kubernetes.RestClientConfig()
			if err != nil {
				return err
			}
			err = execute("POST", req.URL(), config, nil, &stdout, &stderr, tty)

			fmt.Println(strings.TrimSpace(stdout.String()))

			fmt.Printf(`Thanks for using kink.
Pod %s created successfully!
You can view the logs by running the following command:
$ kubectl logs -f %s -n %s 
$ You can start using your internal KinD cluster by running the following command:
$ KUBECONFIG=%s kubectl get nodes -o wide`, outputPath+"/kubeconfig", podName, podName, namespace)
			return nil
		},
	}

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get current directory: %v\n", err)
	}

	cmd.Flags().StringVarP(&k8sVersion, "kubernetes-version", "k", types.ImageTag, "Desired version of Kubernetes")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Target namespace")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", currDir, "Output path for kubeconfig")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "timeout for wait")

	return cmd
}

func execute(method string, url *url.URL, config *rest.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    tty,
	})
}

func ptrbool(p bool) *bool {
	return &p
}

func init() {
	rootCmd.AddCommand(NewCmdRun())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
