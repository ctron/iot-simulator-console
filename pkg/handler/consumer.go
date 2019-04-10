/*******************************************************************************
 * Copyright (c) 2018, 2019 Red Hat Inc
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
	"time"

	"github.com/ctron/iot-simulator-console/pkg/utils"

	"github.com/ctron/iot-simulator-console/pkg/data"

	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isConsumer(obj metav1.Object) bool {
	labels := obj.GetLabels()
	return labels["iot.simulator.message.type"] != "" &&
		labels["iot.simulator.app"] == "consumer"
}

func (c *controller) fillConsumer(tenants *map[string]*data.Tenant, obj metav1.Object, replicas int) {
	if !isConsumer(obj) {
		return
	}

	tenant, component := c.fillCommon(tenants, obj, replicas)
	if tenant == nil {
		return
	}

	// get current mps

	mps, err := c.metricsClient.QuerySingle(
		context.TODO(),
		fmt.Sprintf(`sum(irate(messages_received_total{type="%s",tenant="%s"}[1m]))`, component.Type, tenant.Name),
	)
	if err != nil {
		log.Warn("Failed to query metrics ", err.Error())
	}

	pps, err := c.metricsClient.QuerySingle(
		context.TODO(),
		fmt.Sprintf(`sum(irate(payload_received_total{type="%s",tenant="%s"}[1m]))`, component.Type, tenant.Name),
	)
	if err != nil {
		log.Warn("Failed to query metrics ", err.Error())
	}

	// get last values

	history, err := c.metricsClient.QueryArray(
		context.TODO(), 5*time.Minute, 10*time.Second,
		fmt.Sprintf(`sum(rate(messages_received_total{type="%s",tenant="%s"}[1m]))`, component.Type, tenant.Name),
		[]string{"msg/s"},
	)
	if err != nil {
		log.Warn("Failed to query metrics ", err.Error())
	}

	// assemble

	component.Good = mps != nil && *mps > 0

	tenant.Consumers = append(tenant.Consumers, data.Consumer{
		Component:         component,
		MessagesPerSecond: mps,
		PayloadPerSecond:  utils.FilterNaN(pps),
		MessagesHistory:   history,
	})

}
