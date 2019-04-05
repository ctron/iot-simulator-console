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
	"github.com/ctron/iot-simulator-console/pkg/data"
	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getMessageType(obj metav1.Object) string {
	labels := obj.GetLabels()
	return labels["iot.simulator.message.type"]
}

func (c *controller) fillCommon(tenants *map[string]*data.Tenant, obj metav1.Object, replicas int) (*data.Tenant, data.Component) {

	tenantName := obj.GetLabels()["iot.simulator.tenant"]
	if tenantName == "" {
		return nil, data.Component{}
	}

	tenant := registerTenant(tenants, tenantName)

	messageType := getMessageType(obj)
	if messageType == "" {
		log.Warn("Missing message type")
		return nil, data.Component{}
	}

	return tenant, data.Component{
		Type:     messageType,
		Replicas: uint32(replicas),
	}
}
