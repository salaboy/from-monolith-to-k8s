package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	clientKubernetes "k8s.io/client-go/kubernetes"
	clientCmd "k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	factory := EnvironmentFactory{
		Environment: Environment{
			Name:      "dev-environment",
			Kind:      "dev",
			Size:      "medium",
			GitOpsURL: "",
		},
	}

	factory.createEnvironment()

}

type Environment struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Size      string `json:"size"`
	GitOpsURL string `json:"git_ops_url"`
}

type EnvironmentFactory struct {
	Environment Environment
}

func (envFactory *EnvironmentFactory) createEnvironment() error {
	pulumi.Run(func(ctx *pulumi.Context) error {

		if envFactory.Environment.Kind == "prod" {
			//Create a GKE Cluster in Google Cloud
			// Obtain credentials
			// use gcloud and credentials to connect to the newly created cluster
			// install Helm Chart for conference app
		}

		if envFactory.Environment.Kind == "dev" {
			namespace, err := corev1.NewNamespace(ctx, envFactory.Environment.Name, nil)

			vClusterRelease, err := helm.NewRelease(ctx, envFactory.Environment.Name, &helm.ReleaseArgs{
				Chart:     pulumi.String("vcluster"),
				Namespace: namespace.Metadata.Name(),
				RepositoryOpts: &helm.RepositoryOptsArgs{
					Repo: pulumi.String("https://charts.loft.sh"),
				},
				Version: pulumi.String("0.10.2"),
				Values: pulumi.Map{
					"Syncer": pulumi.Map{
						"extraArgs": pulumi.Map{
							"--out-kube-config-secret": pulumi.String(envFactory.Environment.Name + "-secret"),
						},
					},
				},
			})

			vClusterSecretData := pulumi.All(vClusterRelease.Status.Namespace(), vClusterRelease.Status.Name()).ApplyT(func(r interface{}) interface{} {
				arr := r.([]interface{})
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				vClusterSecret, err := corev1.GetSecret(ctx, "secret", pulumi.ID(fmt.Sprintf("%s/vc-%s", *namespace, *name)), nil)
				if err != nil {
					return ""
				}
				t := vClusterSecret.Data.ApplyT(func(r interface{}) string {
					secretMap := r.(map[string]string)
					payload, err := decodePayload([]byte(secretMap["config"]))
					if err != nil {
						return ""
					}
					s := string(payload)

					log.Println("Config \n" + s)
					return s
				})

				return t
			})

			ctx.Export("vClusterSecretData", vClusterSecretData)

			//client-go port- forward https://stackoverflow.com/questions/59027739/upgrading-connection-error-in-port-forwarding-via-client-go

			portForward(ctx, vClusterRelease, vClusterSecretData)

			if err != nil {
				return err
			}

			return nil
		}
		return nil
	})
	return nil
}

func portForward(ctx *pulumi.Context, vClusterRelease *helm.Release, vClusterSecretData pulumi.Output) {

	stopCh := make(<-chan struct{})
	readyCh := make(chan struct{})

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientCmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	// create the clientset
	clientset, err := clientKubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	pulumi.All(vClusterRelease.Status.Namespace(), vClusterRelease.Status.Name()).ApplyT(func(r interface{}) (interface{}, error) {
		arr := r.([]interface{})
		namespace := arr[0].(*string)
		name := arr[1].(*string)
		log.Println("Namespace/Service: " + fmt.Sprintf("%s/%s", *namespace, *name))

		reqURL := clientset.CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(fmt.Sprintf("%s", *namespace)).
			Name(fmt.Sprintf("%s-0", *name)).
			SubResource("portforward").URL()

		log.Println("URL : " + fmt.Sprintf("%s", reqURL))

		transport, upgrader, err := spdy.RoundTripperFor(config) // I need to get the Kubeconfig for the Host Cluster
		if err != nil {
			log.Println("ERROR in RoundTripperFor")
			log.Fatal(err)
		}
		dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, reqURL)
		fw, err := portforward.New(dialer, []string{"8443:8443"}, stopCh, readyCh, os.Stdout, os.Stdout)
		if err != nil {
			log.Println("ERROR in NewDialer")
			log.Fatal(err)
		}
		if err := fw.ForwardPorts(); err != nil {
			log.Println("ERROR in ForwardPorts")
			log.Fatal(err)
		}

		return nil, nil
	})

	go func() {
		<-readyCh

		args := kubernetes.ProviderArgs{
			Kubeconfig: pulumi.Sprintf("%s", vClusterSecretData),
		}
		provider, err := kubernetes.NewProvider(ctx, "vcluster", &args) //use  vClusterSecretData here
		if err != nil {
			log.Println("ERROR in NewProvider")
			log.Fatal(err)
		}

		installConferenceApplication(ctx, provider)



		<-stopCh
	}()
}
func installConferenceApplication(ctx *pulumi.Context, provider *kubernetes.Provider){
	helm.NewRelease(ctx, "conference", &helm.ReleaseArgs{
		Chart: pulumi.String("fmtok8s-conference-chart"),
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://salaboy.github.io/helm/"),
		},
		Version:   pulumi.String("v0.1.1"),
		Values:    pulumi.Map{},
		SkipAwait: pulumi.BoolPtr(true),
	}, pulumi.Provider(provider))
}

func decodePayload(body []byte) ([]byte, error) {
	//Base64 Decode
	b64 := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	n, err := base64.StdEncoding.Decode(b64, body)
	if err != nil {
		return nil, err
	}
	return b64[:n], nil
}
