// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmccormick2001/rqlite-operator/pkg/apis/rqcluster/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/scheme"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	version = "0.0.4"
)

var schemeGroupVersion = v1alpha1.SchemeGroupVersion

func main() {
	fmt.Printf("my-custom-test version is %s\n", version)

	rqcluster := v1alpha1.RqclusterSpec{}
	rqcluster.Size = 12

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := os.Getenv("POD_NAMESPACE")
	fmt.Printf("POD_NAMESPACE is %s\n", namespace)
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	_, err = clientset.CoreV1().Pods(namespace).Get("example-xxxxx", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("Found pod\n")
	}

	s := runtime.NewScheme()
	v1alpha1.SchemeBuilder.AddToScheme(s)

	crdConfig := *config
	sgv := schema.GroupVersion{Group: "rqcluster.example.com", Version: "v1alpha1"}
	crdConfig.GroupVersion = &sgv
	crdConfig.APIPath = "/apis"
	crdConfig.ContentType = runtime.ContentTypeJSON
	crdConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	crName := os.Getenv("CR_NAME")
	fmt.Printf("CR_NAME is %s\n", crName)

	restclient, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		panic(err.Error())
	}

	result := v1alpha1.Rqcluster{}
	err = restclient.
		Get().
		Namespace(namespace).
		Resource("rqclusters").
		Name(crName).
		Do().
		Into(&result)

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("found the CR %s\n", crName)

	tmp := os.Getenv("EXPECTED_NODES")
	expectedNodes, err := strconv.Atoi(tmp)
	if err != nil {
		fmt.Printf("error in parsing EXPECTED_NODES")
		os.Exit(2)
	}
	actualNodes := len(result.Status.Nodes)
	fmt.Printf("this test expects %d status nodes and got %d\n", expectedNodes, actualNodes)
	if expectedNodes != actualNodes {
		fmt.Printf("this test failed")
		os.Exit(1)
	}
	fmt.Printf("this test passed")

}
