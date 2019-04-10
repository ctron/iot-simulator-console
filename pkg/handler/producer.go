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

package handler

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/ctron/iot-simulator-console/pkg/utils"

	corev1 "k8s.io/api/core/v1"

	"github.com/ctron/iot-simulator-console/pkg/data"
	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isProducer(obj metav1.Object) bool {
	labels := obj.GetLabels()
	return labels["iot.simulator.app"] == "producer"
}

func (c *controller) fillProducer(tenants *map[string]*data.Tenant, obj metav1.Object, pod *corev1.PodTemplateSpec, replicas int) {
	if !isProducer(obj) {
		return
	}

	tenant, component := c.fillCommon(tenants, obj, replicas)

	if tenant == nil {
		return
	}

	protocol := obj.GetLabels()["iot.simulator.producer.protocol"]
	if protocol == "" {
		return
	}

	ctx := context.TODO()

	// messages scheduled/s

	mpsScheduled, err := c.metricsClient.QuerySingle(ctx,
		fmt.Sprintf(`sum(irate(messages_scheduled_total{tenant="%s",type="%s",protocol="%s"}[1m]))`,
			tenant.Name, component.Type, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s scheduled: %v", err)
	}

	// messages sent/s

	mpsSent, err := c.metricsClient.QuerySingle(ctx,
		fmt.Sprintf(`sum(irate(messages_sent_total{type="%s",tenant="%s",protocol="%s"}[1m]))`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s sent: %v", err)
	}

	// failures/s

	mpsFailed, err := c.metricsClient.QuerySingle(ctx,
		fmt.Sprintf(`sum(irate(messages_failure_total{type="%s",tenant="%s",protocol="%s"}[1m]))`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s failed: %v", err)
	}

	// errors/s

	mpsErrored, err := c.metricsClient.QueryMap(ctx,
		fmt.Sprintf(`sum(irate(messages_error_total{type="%s",tenant="%s",protocol="%s"}[1m])) by (code)`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s errored: %v", err)
	}

	// established connections

	conEst, err := c.metricsClient.QuerySingle(ctx,
		fmt.Sprintf(`sum(connections{type="%s",tenant="%s",protocol="%s"})`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query connections established: %v", err)
	}

	// RTT

	rtt, err := c.metricsClient.QuerySingle(ctx,
		fmt.Sprintf(`avg(irate(messages_duration_seconds_sum{type="%[1]s",tenant="%[2]s",protocol="%[3]s"}[1m])/irate(messages_duration_seconds_count{type="%[1]s",tenant="%[2]s",protocol="%[3]s"}[1m])*1000.0)`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query connections established: %v", err)
	}

	// avg(irate(messages_duration_seconds_sum{type="%[1]s",tenant="%[2]s",protocol="%[3]s"}[1m])/irate(messages_duration_seconds_count{type="%[1]s",tenant="%[2]s",protocol="%[3]s"}[1m])*1000.0)

	var chartData []data.ChartEntry
	if mpsSent != nil && mpsErrored != nil {
		chartData = []data.ChartEntry{
			{Key: "Success", Value: *mpsSent},
		}
		keys := make([]string, 0, len(mpsErrored))
		for k := range mpsErrored {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := mpsErrored[k]
			if v > 0 {
				chartData = append(chartData, data.ChartEntry{k, v})
			}
		}
	}

	mpsCfg, conCfg := calcConfiguredMessagesPerSecond(pod, replicas)

	switch protocol {
	case "mqtt":
		component.Good = conCfg != nil && conEst != nil && *conCfg == *conEst
	default:
		component.Good = isGoodHttp(mpsScheduled, mpsSent, mpsFailed)
	}

	tenant.Producers = append(tenant.Producers, data.Producer{
		Component: component,
		Protocol:  protocol,

		MessagesPerSecondConfigured: mpsCfg,
		MessagesPerSecondScheduled:  mpsScheduled,
		MessagesPerSecondSent:       mpsSent,
		MessagesPerSecondFailed:     mpsFailed,
		MessagesPerSecondErrored:    sum(mpsErrored),

		RoundTripTime: utils.FilterNaN(rtt),

		ConnectionsConfigured:  conCfg,
		ConnectionsEstablished: conEst,

		ChartData:   chartData,
		ChartLegend: makeLegend(chartData),
	})

}

func isGoodHttp(mpsScheduled, mpsSent, mpsFailed *float64) bool {

	if mpsScheduled == nil || mpsSent == nil || mpsFailed == nil {
		return false
	}

	if *mpsSent < *mpsScheduled*0.95 {
		// allow 5% miss
		return false
	}

	if *mpsFailed > *mpsSent*0.05 {
		// allow 5% failure
		return false
	}

	return true
}

func calcConfiguredMessagesPerSecond(pod *corev1.PodTemplateSpec, replicas int) (mpsConfigured *float64, connectionsConfigured *float64) {

	var numDevices *float64
	var period *float64

	for _, c := range pod.Spec.Containers {

		for _, e := range c.Env {
			switch e.Name {
			case "PERIOD_MS":
				if v, err := strconv.ParseFloat(e.Value, 64); err == nil {
					if v > 0.0 {
						period = &v
					}
				}
			case "NUM_DEVICES":
				if v, err := strconv.ParseFloat(e.Value, 64); err == nil {
					numDevices = &v
				}
			default:
			}
		}
	}

	if numDevices != nil && period != nil {
		cons := (*numDevices) * float64(replicas)
		msgs := (*numDevices) * (1000.0 / (*period)) * float64(replicas)
		if math.IsNaN(msgs) {
			return nil, nil
		}
		return &msgs, &cons
	}

	return nil, nil

}

func sum(data map[string]float64) *float64 {
	if data == nil {
		return nil
	}

	var result float64 = 0.0

	for _, v := range data {
		result += v
	}

	return &result
}

func makeLegend(d []data.ChartEntry) []data.ChartLegendEntry {
	if d == nil {
		return nil
	}

	result := make([]data.ChartLegendEntry, len(d))

	for i, e := range d {
		result[i] = data.ChartLegendEntry{Name: e.Key}
	}

	return result
}
