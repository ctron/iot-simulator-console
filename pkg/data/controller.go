/*******************************************************************************
 * Copyright (c) 2018 Red Hat Inc
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 *******************************************************************************/

package data

import (
	"github.com/ctron/iot-simulator-console/pkg/metrics"
	"github.com/ctron/operator-tools/pkg/install/openshift"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type controller struct {
	namespace     string
	client        *kubernetes.Clientset
	appsclient    *appsv1.AppsV1Client
	metricsClient *metrics.MetricsClient
}

func NewController(namespace string, client *kubernetes.Clientset, appsclient *appsv1.AppsV1Client, metricsClient *metrics.MetricsClient) *controller {
	return &controller{
		namespace:     namespace,
		client:        client,
		appsclient:    appsclient,
		metricsClient: metricsClient,
	}
}

func registerTenant(tenants *map[string]*Tenant, tenantName string) *Tenant {
	tenant, ok := (*tenants)[tenantName]
	if !ok {
		tenant = &Tenant{Name: tenantName}
		(*tenants)[tenantName] = tenant
	}
	return tenant
}

type OverviewProcessor func(obj metav1.Object, pod *corev1.PodTemplateSpec, replicas int)

func (c *controller) ProcessDeploymentConfigs(processor OverviewProcessor) error {

	items, err := c.appsclient.DeploymentConfigs(c.namespace).
		List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, i := range items.Items {
		processor(&i, i.Spec.Template, int(i.Spec.Replicas))
	}

	return nil
}

func (c *controller) ProcessDeployments(processor OverviewProcessor) error {

	items, err := c.client.AppsV1().Deployments(c.namespace).
		List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, i := range items.Items {
		var r int
		if i.Spec.Replicas != nil {
			r = int(*i.Spec.Replicas)
		} else {
			r = 1
		}
		processor(&i, &i.Spec.Template, r)
	}

	return nil
}

func (c *controller) BuildOverview() (*Overview, error) {

	tenants := map[string]*Tenant{}

	fn := func(obj metav1.Object, pod *corev1.PodTemplateSpec, replicas int) {
		c.fillConsumer(&tenants, obj, replicas)
		c.fillProducer(&tenants, obj, pod, replicas)
	}

	if openshift.IsOpenshift() {
		if err := c.ProcessDeploymentConfigs(fn); err != nil {
			return nil, err
		}
	} else {
		if err := c.ProcessDeployments(fn); err != nil {
			return nil, err
		}
	}

	return &Overview{
		Tenants: makeTenants(tenants),
	}, nil
}

func makeTenants(t map[string]*Tenant) []Tenant {
	var result = make([]Tenant, 0, len(t))
	for _, v := range t {
		result = append(result, *v)
	}
	return result
}
