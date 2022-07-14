/*
Copyright 2022.

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

package controllers

import (
	"context"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	conferencev1 "github.com/salaboy/conference-controller/api/v1"
)

// ConferenceReconciler reconciles a Conference object
type ConferenceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=kubebuilder.conference.salaboy.com,resources=conferences,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubebuilder.conference.salaboy.com,resources=conferences/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubebuilder.conference.salaboy.com,resources=conferences/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Conference object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *ConferenceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	log.V(1).Info("Conference Controller:", "Namespace del recurso", req.Namespace, "Nombre del recurso", req.Name)

	var conference conferencev1.Conference
	if err := r.Get(ctx, req.NamespacedName, &conference); err != nil {
		log.Error(err, "unable to fetch Conference")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var services v1.ServiceList
	if err := r.List(ctx, &services, client.InNamespace(conference.Spec.Namespace), client.MatchingLabels{"draft": "draft-app"}); err != nil {
		log.Error(err, "unable to list conference services")
		return ctrl.Result{}, err
	}

	for _, service := range services.Items {
		log.V(1).Info("Services", "Service", service.Name)
		//Call each service and do a test to see if it is operational
		resp, err := http.Get("http://" + service.Name + "." + conference.Spec.Namespace + ".svc.cluster.local/info")
		if err != nil {
			if service.Name == "fmtok8s-frontend" {
				conference.Status.FrontendReady = false
			}
			if service.Name == "fmtok8s-agenda" {
				conference.Status.AgendaServiceReady = false
			}
			if service.Name == "fmtok8s-c4p" {
				conference.Status.C4pServiceReady = false
			}
			if service.Name == "fmtok8s-email" {
				conference.Status.EmailServiceReady = false
			}
			log.V(1).Error(err, "ServicesCheckNotOK", "Service", service.Name)
			continue
		} else {
			log.V(1).Info("ServicesCheckOK", "Service", service.Name, "Response Code", resp.StatusCode)
			if resp.StatusCode == 200 {
				if service.Name == "fmtok8s-frontend" {
					conference.Status.FrontendReady = true
				}
				if service.Name == "fmtok8s-agenda" {
					conference.Status.AgendaServiceReady = true
				}
				if service.Name == "fmtok8s-c4p" {
					conference.Status.C4pServiceReady = true
				}
				if service.Name == "fmtok8s-email" {
					conference.Status.EmailServiceReady = true
				}
			}
		}

		if service.Name == "fmtok8s-frontend" {
			conference.Status.URL = service.Status.LoadBalancer.Ingress[0].IP
			log.V(1).Info("Frontend", "IP", conference.Status.URL)
		}

	}
	if conference.Status.AgendaServiceReady == true && conference.Status.EmailServiceReady == true &&
		conference.Status.C4pServiceReady == true && conference.Status.FrontendReady == true {
		conference.Status.Ready = true

		if conference.Spec.ProductionTestEnabled {
			//       - Check if exists, if not create and trigger and report to status
			var deployment appsV1.Deployment
			var deploymentFound bool
			if err := r.Get(ctx, client.ObjectKey{
				Namespace: req.Namespace,
				Name:      "kubebuilder-production-test",
			}, &deployment); err != nil {
				deploymentFound = false
			} else {
				deploymentFound = true
			}

			if !deploymentFound {
				deployment, _ := createProductionTestDeployment(&conference, req.Namespace, r.Scheme)
				if err := r.Create(ctx, &deployment); err != nil {
					log.Error(err, "unable create production check deployment")
					return ctrl.Result{}, err
				}

				// @TODO: call production test endpoint and report back to status
				conference.Status.ProdTests = true
			}
		}

	} else {
		conference.Status.Ready = false
	}

	if err := r.Status().Update(ctx, &conference); err != nil {
		log.Error(err, "unable to update Conference status")
		return ctrl.Result{}, err
	}

	requeue := ctrl.Result{RequeueAfter: time.Second * 5}

	return requeue, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConferenceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&conferencev1.Conference{}).
		Complete(r)
}

func createProductionTestDeployment(conference *conferencev1.Conference, namespace string, runtimeScheme *runtime.Scheme) (appsV1.Deployment, error) {

	replicas := int32(1)
	deployment := appsV1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubebuilder-production-test",
			Namespace: namespace,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"label": "labelvalue",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"label": "labelvalue",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "production-test",
							Image: "salaboy/conference-production-test:0.1.0",
							Env: []v1.EnvVar{
								{
									Name:  "testvar",
									Value: "testvarvalue",
								},
							},
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(conference, &deployment, runtimeScheme); err != nil {
		return appsV1.Deployment{}, err
	}
	return deployment, nil
}
