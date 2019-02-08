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
)

func isConsumer(dc *v1.DeploymentConfig) bool {
	return dc.Labels["iot.simulator.consume.type"] != ""
}

func (c *controller) fillConsumer(tenants *map[string]*Tenant, dc *v1.DeploymentConfig) {
	if !isConsumer(dc) {
		return
	}

	tenant, component := c.fillCommon(tenants, dc)
	if tenant == nil {
		return
	}

	mps, err := c.metricsClient.QuerySingle(
		context.TODO(),
		fmt.Sprintf(`sum(irate(messages_received_total{type="%s",tenant="%s"}[1m]))`, component.Type, tenant.Name),
	)
	if err != nil {
		log.Warn("Failed to query metrics", err.Error())
	}

	tenant.Consumers = append(tenant.Consumers, Consumer{
		Component:         component,
		MessagesPerSecond: mps,
	})

}
