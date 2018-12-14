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

type Overview struct {
	Tenants []Tenant `json:"tenants"`
}

type Tenant struct {
	Name       string           `json:"name"`
	Statistics TenantStatistics `json:"statistics"`
}

type TypeStatistics struct {
	Producer ProducerRate `json:"producer""`
	Consumer ConsumerRate `json:"consumer"`
}

type ConsumerRate struct {
	Replicas     uint32  `json:"replicas"`
	ReceivedRate float64 `json:"receivedRate"`
}

type ProducerRate struct {
	Replicas          uint32  `json:"replicas"`
	MessagesPerSecond float64 `json:"messagesPerSecond"`
}

type TenantStatistics struct {
	Telemetry TypeStatistics `json:"telemetry"`
	Event     TypeStatistics `json:"event"`
}
