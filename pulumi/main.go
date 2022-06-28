package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Get the Pulumi API token.
		c := config.New(ctx, "")
		pulumiAccessToken := c.Require("pulumiAccessToken")

		// Create the API token as a Kubernetes Secret.
		accessToken, err := corev1.NewSecret(ctx, "accesstoken", &corev1.SecretArgs{
			StringData: pulumi.StringMap{"accessToken": pulumi.String(pulumiAccessToken)},
		})
		if err != nil {
			return err
		}

		// Create an NGINX deployment in-cluster.
		_, err = apiextensions.NewCustomResource(ctx, "dev", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("pulumi.com/v1"),
			Kind:       pulumi.String("Stack"),
			OtherFields: kubernetes.UntypedArgs{
				"spec": map[string]interface{}{
					"envRefs": pulumi.Map{
						"PULUMI_ACCESS_TOKEN": pulumi.Map{
							"type": pulumi.String("Secret"),
							"secret": pulumi.Map{
								"name": accessToken.Metadata.Name(),
								"key":  pulumi.String("accessToken"),
							},
						},
					},
					"stack":             "dev",
					"projectRepo":       "https://github.com/salaboy/from-monolith-to-k8s/",
					//"commit":            "2b0889718d3e63feeb6079ccd5e4488d8601e353",
					"branch":         "refs/heads/main", // Alternatively, track master branch.
					"repoDir": "pulumi/helm/",
					"destroyOnFinalize": true,
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{accessToken}))
		return err
	})

}
