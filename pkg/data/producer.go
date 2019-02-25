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
	"github.com/prometheus/common/log"
	"math"
	"strconv"
)

func isProducer(dc *v1.DeploymentConfig) bool {
	return dc.Labels["iot.simulator.app"] == "producer"
}

func calcConfiguredMessagesPerSecond(dc *v1.DeploymentConfig) *float64 {

	var numDevices *float64
	var period *float64

	for _, c := range dc.Spec.Template.Spec.Containers {

		for _, e := range c.Env {
			switch e.Name {
			case "TELEMETRY_MS":
				fallthrough
			case "EVENT_MS":
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
		f := (*numDevices) * (1000.0 / (*period)) * float64(dc.Spec.Replicas)
		if math.IsNaN(f) {
			return nil
		}
		return &f
	}

	return nil

}

func (c *controller) fillProducer(tenants *map[string]*Tenant, dc *v1.DeploymentConfig) {
	if !isProducer(dc) {
		return
	}

	tenant, component := c.fillCommon(tenants, dc)

	if tenant == nil {
		return
	}

	protocol := dc.Labels["iot.simulator.producer.protocol"]
	if protocol == "" {
		return
	}

	mpsScheduled, err := c.metricsClient.QuerySingle(context.TODO(),
		fmt.Sprintf(`sum(irate(messages_scheduled_total{tenant="%s",type="%s",protocol="%s"}[1m]))`,
			tenant.Name, component.Type, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s scheduled: %v", err)
	}

	mpsSent, err := c.metricsClient.QuerySingle(context.TODO(),
		fmt.Sprintf(`sum(irate(messages_success_total{type="%s",tenant="%s",protocol="%s"}[1m]))`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s sent: %v", err)
	}

	mpsFailed, err := c.metricsClient.QuerySingle(context.TODO(),
		fmt.Sprintf(`sum(irate(messages_failure_total{type="%s",tenant="%s",protocol="%s"}[1m]))`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s failed: %v", err)
	}
	mpsErrored, err := c.metricsClient.QuerySingle(context.TODO(),
		fmt.Sprintf(`sum(irate(messages_error_total{type="%s",tenant="%s",protocol="%s"}[1m]))`,
			component.Type, tenant.Name, protocol))

	if err != nil {
		log.Warnf("Failed to query msg/s errored: %v", err)
	}

	tenant.Producers = append(tenant.Producers, Producer{
		Component:                   component,
		Protocol:                    protocol,
		MessagesPerSecondConfigured: calcConfiguredMessagesPerSecond(dc),
		MessagesPerSecondScheduled:  mpsScheduled,
		MessagesPerSecondSent:       mpsSent,
		MessagesPerSecondFailed:     mpsFailed,
		MessagesPerSecondErrored:    mpsErrored,
	})
}
