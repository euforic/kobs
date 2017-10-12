package kobs

import (
	"log"
	"reflect"
	"time"

	"github.com/segmentio/ksuid"
	batchv1 "k8s.io/api/batch/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Manager struct
type Manager struct {
	Client      *kubernetes.Clientset
	podInformer cache.SharedIndexInformer
	stopCh      chan struct{}
}

// New creates a new Manager and returns a pointer to it
func New(client *kubernetes.Clientset) *Manager {
	// use in cluster config if no client provided
	if client == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		client = clientset
	}

	mgr := Manager{
		Client: client,
		stopCh: make(chan struct{}),
	}

	// Create informer for watching Namespaces
	mgr.podInformer = cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.BatchV1().Jobs("").List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.BatchV1().Jobs("").Watch(options)
			},
		},
		&batchv1.Job{},
		time.Second*30,
		cache.Indexers{},
	)

	mgr.podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, cur interface{}) {
			o := old.(*batchv1.Job)
			c := cur.(*batchv1.Job)
			if !reflect.DeepEqual(o, c) {
				if c.Status.Active == 0 {
					mgr.Delete(c.Name, c.Namespace)
				}
			}
		},
	})

	return &mgr
}

// Create will create a new k8 job in the cluster
func (m *Manager) Create(j *batchv1.Job) (*batchv1.Job, error) {
	if j.Labels == nil {
		j.Labels = map[string]string{}
	}

	if j.Namespace == "" {
		j.Namespace = "default"
	}

	j.Labels["kobs-id"] = ksuid.New().String()
	job, err := m.Client.BatchV1().Jobs(j.Namespace).Create(j)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Get will get a job's details in the k8 cluster
func (m *Manager) Get(name string, namespace string) (*batchv1.Job, error) {
	job, err := m.Client.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Update will update a job in the k8 cluster
func (m *Manager) Update(j *batchv1.Job) (*batchv1.Job, error) {
	job, err := m.Client.BatchV1().Jobs(j.Namespace).Update(j)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Delete will delete a job in the k8 cluster
func (m *Manager) Delete(name string, namespace string) error {
	delBG := metav1.DeletePropagationBackground

	if err := m.Client.BatchV1().Jobs(namespace).Delete(
		name, &metav1.DeleteOptions{PropagationPolicy: &delBG}); err != nil {
		return err
	}
	return nil
}

// List will list all the jobs in the k8 cluster
func (m *Manager) List(namespace string) (*batchv1.JobList, error) {
	jobs, err := m.Client.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// Start starts the process for listening for job changes and acting upon those changes
func (m *Manager) Start() {
	log.Printf("Listening for changes...")
	m.podInformer.Run(m.stopCh)
}

// Stop stops the process for listening for job changes
func (m *Manager) Stop() {
	close(m.stopCh)
	log.Println("Stopped listening for changes.")
}
