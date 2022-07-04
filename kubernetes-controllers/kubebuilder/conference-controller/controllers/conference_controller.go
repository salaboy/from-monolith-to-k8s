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
	v1 "k8s.io/api/core/v1"
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

//+kubebuilder:rbac:groups=conference.salaboy.com,resources=conferences,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=conference.salaboy.com,resources=conferences/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=conference.salaboy.com,resources=conferences/finalizers,verbs=update

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
	if err := r.List(ctx, &services, client.InNamespace(req.Namespace), client.MatchingLabels{"draft": "draft-app"}); err != nil {
		log.Error(err, "unable to list conference services")
		return ctrl.Result{}, err
	}

	for _, service := range services.Items {
		log.V(1).Info("Services", "Service", service.Name)
		//Call each service and do a test to see if it is operational
		if service.Name == "fmtok8s-frontend" {
			conference.Status.URL = service.Status.LoadBalancer.Ingress[0].IP
			log.V(1).Info("Frontend", "IP", conference.Status.URL)
		}
	}
	conference.Status.Ready = true

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
