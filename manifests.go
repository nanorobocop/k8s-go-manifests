package main

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	replicas      = int32(3)
	quantity1Gi   = resource.MustParse("1Gi")
	quantity100Mi = resource.MustParse("100Mi")
	quantity1     = resource.MustParse("1")
	quantity100m  = resource.MustParse("100m")
	port          = intstr.FromInt(8080)

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{"app": "nginx"},
			Name:      "nginx",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "VAR",
									Value: "abc",
								},
								{
									Name: "VAR_CONFIG",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "nginx-config",
											},
											Key: "a",
										},
									},
								},
								{
									Name: "SECRET_VAR",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "nginx-secret",
											},
											Key: "secret",
										},
									},
								},
							},
							Image: "nginx",
							Name:  "nginx",
							Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    quantity1,
									corev1.ResourceMemory: quantity1Gi,
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    quantity100m,
									corev1.ResourceMemory: quantity100Mi,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/",
									SubPath:   "nginx.conf",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{
									Name: "nginx-config"}},
							},
						},
					},
				},
			},
		},
	}

	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{"app": "nginx"},
			Name:      "nginx",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "nginx",
					Port:       80,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: port,
				},
			},
			Selector: map[string]string{"app": "nginx"},
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}

	nginxConf = `worker_processes  1;
events {
    worker_connections  1024;
}
http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;
    server {
        listen       8080;
        server_name  localhost;
        location / {
            root   html;
            index  index.html index.htm;
        }
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}`

	configMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{"app": "nginx"},
			Name:      "nginx-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"nginx.conf": nginxConf,
		},
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{"app": "nginx"},
			Name:      "nginx-secret",
			Namespace: "default",
		},
		StringData: map[string]string{
			"secret": "topSecret123",
		},
	}
)
