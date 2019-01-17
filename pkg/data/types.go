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
	Name      string     `json:"name"`
	Consumers []Consumer `json:"consumers"`
}

type Consumer struct {
	Type              string   `json:"type"`
	Replicas          uint32   `json:"replicas"`
	MessagesPerSecond *float64 `json:"messagesPerSecond"`
}
