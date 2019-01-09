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
	"github.com/openshift/api/apps/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"

	// "github.com/prometheus/client_golang/api"
	// "github.com/prometheus/common/model"
)

type controller struct {
	namespace  string
	client     *kubernetes.Clientset
	appsclient *appsv1.AppsV1Client
}

func NewController(namespace string, client *kubernetes.Clientset, appsclient *appsv1.AppsV1Client) *controller {
	return &controller{
		namespace:  namespace,
		client:     client,
		appsclient: appsclient,
	}
}

func isConsumer(labels *map[string]string) bool {
	return (*labels)["iot.simulator.app"] == "consumer"
}
func isProducer(labels *map[string]string) bool {
	return (*labels)["iot.simulator.app"] == "producer"
}

func (c *controller) BuildOverview() (*Overview, error) {

	items, err := c.appsclient.DeploymentConfigs(c.namespace).
		List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	tenants := map[string]*Tenant{}

	for _, i := range items.Items {
		if val, ok := i.Labels["iot.simulator.tenant"]; ok {

			tenant, ok := tenants[val]
			if !ok {
				tenant = &Tenant{Name: val}
				tenants[val] = tenant
			}

			var stats *TypeStatistics
			switch i.Labels["iot.simulator.type"] {
			case "telemetry":
				stats = &tenant.Statistics.Telemetry
			case "event":
				stats = &tenant.Statistics.Event
			default:
				continue
			}

			if isConsumer(&i.Labels) {
				stats.Consumer.Replicas += uint32(i.Spec.Replicas)
			}
			if isProducer(&i.Labels) {
				fillForProducer(stats, &i)
			}

		}
	}

	return &Overview{
		Tenants: makeTenants(tenants),
	}, nil
}

func fillForProducer(stats *TypeStatistics, dc *v1.DeploymentConfig) {
	stats.Producer.Replicas += uint32(dc.Spec.Replicas)

	var sendPeriod *int64
	var numDevices *int64

	for _, c := range dc.Spec.Template.Spec.Containers {

		for _, ev := range c.Env {
			switch ev.Name {
			case "TELEMETRY_MS", "PERIOD_MS":
				if i, err := strconv.ParseInt(ev.Value, 10, 32); err == nil {
					sendPeriod = &i
				}
			case "NUM_DEVICES":
				if i, err := strconv.ParseInt(ev.Value, 10, 32); err == nil {
					numDevices = &i
				}
			}
		}
	}

	if sendPeriod != nil && numDevices != nil {
		mps := 1000 / (*sendPeriod) * (*numDevices) * int64(stats.Producer.Replicas)
		stats.Producer.MessagesPerSecond += float64(mps)
	}

}

func makeTenants(t map[string]*Tenant) []Tenant {
	var result = make([]Tenant, 0, len(t))
	for _, v := range t {
		result = append(result, *v)
	}
	return result
}
