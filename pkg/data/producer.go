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
)

func isProducer(dc *v1.DeploymentConfig) bool {
	return dc.Labels["iot.producer"] != ""
}

func (c *controller) fillProducer(tenants *map[string]*Tenant, dc *v1.DeploymentConfig) {
	if !isProducer(dc) {
		return
	}

	tenant, component := c.fillCommon(tenants, dc)

	if tenant == nil {
		return
	}

	tenant.Producers = append(tenant.Producers, Producer{
		Component: component,
		Protocol:  dc.Labels["iot.producer"],
	})
}
