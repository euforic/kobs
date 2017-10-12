package kobs

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	K8Job *batchv1.Job
}

func NewJob(name string, image string, commands ...string) *Job {
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: name,
			Namespace: "default",
			// Labels:      map[string]string{},
			// Annotations: map[string]string{},
		},
		Spec: batchv1.JobSpec{
			Parallelism:           nil, // int32
			Completions:           nil, // int32
			ActiveDeadlineSeconds: nil, // int64
			BackoffLimit:          nil, // int32
			Selector:              &metav1.LabelSelector{
			// MatchLabels:      map[string]string{},
			// MatchExpressions: []metav1.LabelSelectorRequirement{},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
					// Labels:      map[string]string{},
					// Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:    name,
							Image:   image,
							Command: commands,
							Args:    []string{},
						},
					},
					// Volumes: []corev1.Volume{},
					// InitContainers: []corev1.Container{},
					// NodeSelector: "",
					ActiveDeadlineSeconds: nil, // int64
					RestartPolicy:         corev1.RestartPolicyNever,
				},
			},
		},
	}
	return &Job{K8Job: &job}
}
