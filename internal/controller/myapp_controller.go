/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	secondlog "log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/my-org/my-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyAppReconciler reconciles a MyApp object
type MyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// int32Ptr is a helper function to get a pointer to an int32 value
func int32Ptr(i int32) *int32 {
	return &i
}

//+kubebuilder:rbac:groups=apps.example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.example.com,resources=myapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.example.com,resources=myapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Load the kubeconfig file to create a client
	kubeconfig := "/Users/cniackz/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		secondlog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Create a new Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		secondlog.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	// ... START OF ADDED LOGIC:

	log := log.FromContext(ctx)

	// Fetch the MyApp instance
	var myapp appsv1alpha1.MyApp
	if err := r.Get(ctx, req.NamespacedName, &myapp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Print the message from the Custom Resource
	log.Info("Received MyResource", "message", myapp.Spec.Message)

	if "create" == myapp.Spec.Message {
		log.Info("Yes, the message is create, hence we are going to create a statefulset")

		// Define the StatefulSet
		statefulSet := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example-statefulset",
				Namespace: "default",
			},
			Spec: appsv1.StatefulSetSpec{
				ServiceName: "example-service",
				Replicas:    int32Ptr(3),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "example",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "example",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "example-container",
								Image: "nginx:latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "example-pvc",
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Resources: corev1.VolumeResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
			},
		}

		// Create the StatefulSet
		statefulSetsClient := clientset.AppsV1().StatefulSets("default")
		result, err := statefulSetsClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
		if err != nil {
			secondlog.Fatalf("Error creating StatefulSet: %s", err.Error())
		}

		fmt.Printf("Created StatefulSet %s\n", result.GetName())

	}

	// ... END OF ADDED LOGIC

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.MyApp{}).
		Complete(r)
}
