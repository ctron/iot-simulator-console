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

package data

import (
	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getMessageType(obj metav1.Object) string {
	labels := obj.GetLabels()
	return labels["iot.simulator.message.type"]
}

func (c *controller) fillCommon(tenants *map[string]*Tenant, obj metav1.Object, replicas int) (*Tenant, Component) {

	tenantName := obj.GetLabels()["iot.simulator.tenant"]
	if tenantName == "" {
		return nil, Component{}
	}

	tenant := registerTenant(tenants, tenantName)

	messageType := getMessageType(obj)
	if messageType == "" {
		log.Warn("Missing message type")
		return nil, Component{}
	}

	return tenant, Component{
		Type:     messageType,
		Replicas: uint32(replicas),
	}
}
