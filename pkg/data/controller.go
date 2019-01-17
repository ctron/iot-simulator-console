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
	"context"
	"fmt"
	"github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"

	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"
)

type controller struct {
	namespace  string
	client     *kubernetes.Clientset
	appsclient *appsv1.AppsV1Client
	promapi    promapi.API
}

func NewController(namespace string, client *kubernetes.Clientset, appsclient *appsv1.AppsV1Client, promapi promapi.API) *controller {
	return &controller{
		namespace:  namespace,
		client:     client,
		appsclient: appsclient,
		promapi:    promapi,
	}
}

/*
func isConsumer(labels *map[string]string) bool {
	return (*labels)["iot.simulator.app"] == "consumer"
}
func isProducer(labels *map[string]string) bool {
	return (*labels)["iot.simulator.app"] == "producer"
}
*/

func registerTenant(tenants *map[string]*Tenant, tenantName string) *Tenant {
	tenant, ok := (*tenants)[tenantName]
	if !ok {
		tenant = &Tenant{Name: tenantName}
		(*tenants)[tenantName] = tenant
	}
	return tenant
}

func isConsumer(dc *v1.DeploymentConfig) bool {
	return dc.Labels["deploymentconfig"] == "simulator-consumer"
}

func (c *controller) promquery(query string) (*float64, error) {

	s := time.Now()
	e := s.Add(-time.Minute)

	val, err := c.promapi.Query(context.TODO(), query, e)
	if err != nil {
		return nil, err
	}

	switch v := val.(type) {
	case *prommodel.Scalar:
		f := float64(v.Value)
		return &f, nil
	case prommodel.Vector:
		if len(v) > 0 {
			f := float64(v[0].Value)
			return &f, nil
		} else {
			return nil, fmt.Errorf("missing values in vector result")
		}
	default:
		return nil, fmt.Errorf("unknown result type: %v / %v", val.Type().String(), val.String())
	}

}

func (c *controller) fillConsumer(tenants *map[string]*Tenant, dc *v1.DeploymentConfig) {
	if !isConsumer(dc) {
		return
	}

	tenantName, ok := dc.Spec.Template.Labels["iot.simulator.tenant"]
	if !ok {
		return
	}

	tenant := registerTenant(tenants, tenantName)

	mps, err := c.promquery(fmt.Sprintf(`sum(irate(messages_received_total{type="%s",tenant="%s"}[1m]))`, "telemetry", tenantName))
	if err != nil {
		log.Warn("Failed to query metrics", err.Error())
	}

	tenant.Consumers = append(tenant.Consumers, Consumer{
		Type:              "telemetry",
		Replicas:          uint32(dc.Spec.Replicas),
		MessagesPerSecond: mps,
	})
}

func (c *controller) BuildOverview() (*Overview, error) {

	items, err := c.appsclient.DeploymentConfigs(c.namespace).
		List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	tenants := map[string]*Tenant{}

	for _, i := range items.Items {

		c.fillConsumer(&tenants, &i)

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
