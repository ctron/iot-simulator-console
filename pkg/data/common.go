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
	"github.com/prometheus/common/log"
)

func getMessageType(dc *v1.DeploymentConfig) string {
	return dc.Labels["iot.simulator.message.type"]
}

func (c *controller) fillCommon(tenants *map[string]*Tenant, dc *v1.DeploymentConfig) (*Tenant, Component) {

	tenantName := dc.Spec.Template.Labels["iot.simulator.tenant"]
	if tenantName == "" {
		return nil, Component{}
	}

	tenant := registerTenant(tenants, tenantName)

	messageType := getMessageType(dc)
	if messageType == "" {
		log.Warn("Missing message type")
		return nil, Component{}
	}

	return tenant, Component{
		Type:     messageType,
		Replicas: uint32(dc.Spec.Replicas),
	}
}
