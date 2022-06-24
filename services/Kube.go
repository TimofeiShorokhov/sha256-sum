package services

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sha256-sum/models"
	"time"
)

func (s *HashService) ConnToPod() *kubernetes.Clientset {
	log.Printf("### 🚀 PodKicker %s starting...", "0.0.0")

	// Connect to Kubernetes API
	log.Printf("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("### 💻 Connecting to Kubernetes API, using host: %s", config.Host)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func (s *HashService) Podkicker(code int, path string) {
	clientset := s.ConnToPod()
	targetName := os.Getenv("POD_NAME")
	namespace := os.Getenv("NAMESPACE")
	data, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), targetName, metav1.GetOptions{})
	var result models.PodInfo

	if err != nil {
		log.Fatalln(err)
	}
	if code == 1 {
		isRestarting := false

		if isRestarting {
			return
		}

		log.Print("### ⛔ Detected file change")

		patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))

		var err1 error
		_, err1 = clientset.AppsV1().Deployments(namespace).Patch(context.Background(), targetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
		if err1 != nil {
			log.Printf("### 👎 Warning: Failed to patch %s, restart failed: %v", "deployment", err1)
		} else {
			isRestarting = true
			log.Printf("### ✅ Target %s, named %s was restarted!", "deployment", targetName)
		}
	}
	if code == 0 {
		result.PodName = targetName
		result.CreationTime = data.CreationTimestamp.GoString()
		for _, v := range data.Spec.Template.Spec.Containers {
			result.ImageName = v.Image
			result.ContainerName = v.Name
		}
		s.SavingData(s.CheckSum(path), result)
		log.Printf("### ✅ Data was inserted, pod name: %s", targetName)
	}
}

func (s *HashService) Operations(code int, path string) {
	switch {
	case code == 0:
		s.Podkicker(0, path)
	case code == 1:
		check, err := s.GetChangedData(path)
		if err != nil {
			log.Fatalln(err)
		}
		if check == 1 {
			s.repo.Truncate()
			s.Podkicker(1, path)
			fmt.Println("database has changes, truncate successful")
		}
	}
}
