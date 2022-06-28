package main

import (
	"fmt"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		namespace, err := corev1.NewNamespace(ctx, "pulumi-conference", nil)

		conferenceApp, err := helm.NewRelease(ctx, "conference-dev", &helm.ReleaseArgs{
			Version: pulumi.String("0.1.0"),
			Chart:   pulumi.String("fmtok8s-conference-chart"),
			Namespace: namespace.Metadata.Name(),
			RepositoryOpts: &helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://salaboy.github.io/helm/"),
			},
		})

		// Export the ingress IP for Frontend Service frontend.
		frontendIp := pulumi.All(conferenceApp.Status.Namespace(), conferenceApp.Status.Name()).ApplyT(func(r interface{})(interface{}, error){
			arr := r.([]interface{})
			namespace := arr[0].(*string)
			name := arr[1].(*string)
			svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s-frontend", *namespace, *name)), nil)
			if err != nil {
				return "", nil
			}
			return svc.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip(), nil

		})
		ctx.Export("frontendIp", frontendIp)

		if err != nil {
			return err
		}

		return nil
	})

}
