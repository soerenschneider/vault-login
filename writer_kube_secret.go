package main

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const defaultSecretkey = "vault-token"

type KubernetesSecretWriter struct {
	namespace  string
	secretName string
	secretKey  string

	client *kubernetes.Clientset
}

func NewKubernetesSecretWriter(cfg Config) (*KubernetesSecretWriter, error) {
	if cfg.OutputSecretName == "" {
		return nil, errors.New("no secret name supplied")
	}

	if cfg.OutputSecretNamespace == "" {
		return nil, errors.New("no secret namespace supplied")
	}

	if cfg.OutputSecretKey == "" {
		cfg.OutputSecretKey = defaultSecretkey
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubernetesSecretWriter{
		client:     clientset,
		secretName: cfg.OutputSecretName,
		namespace:  cfg.OutputSecretNamespace,
		secretKey:  cfg.OutputSecretKey,
	}, nil
}

func (w *KubernetesSecretWriter) Write(ctx context.Context, data []byte) error {
	secretData := map[string][]byte{
		w.secretKey: data,
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: w.secretName,
		},
		Type: v1.SecretTypeOpaque,
		Data: secretData,
	}

	_, err := w.client.CoreV1().Secrets(w.namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		if kubeErrors.IsNotFound(err) {
			_, err := w.client.CoreV1().Secrets(w.namespace).Create(ctx, secret, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create secret: %v", err)
			}
		} else {
			return fmt.Errorf("failed to write secret: %v", err)
		}
	}

	return err
}
