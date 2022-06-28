package main

import (
	"fmt"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		wordpress, err := helm.NewRelease(ctx, "conference-dev", &helm.ReleaseArgs{
			Version: pulumi.String("0.1.0"),
			Chart:   pulumi.String("fmtok8s-conference-chart"),
			RepositoryOpts: &helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://github.com/salaboy/helm"),
			},
		})

		// Export the ingress IP for Wordpress frontend.
		frontendIp := pulumi.All(wordpress.Status.Namespace(), wordpress.Status.Name()).ApplyT(func(r interface{})(interface{}, error){
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
