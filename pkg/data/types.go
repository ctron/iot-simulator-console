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
	Producers []Producer `json:"producers"`
}

type Component struct {
	Type     string `json:"type"`
	Replicas uint32 `json:"replicas"`
	Good     bool   `json:"good"`
}

type Consumer struct {
	Component `json:",inline"`

	MessagesPerSecond *float64 `json:"messagesPerSecond"`
	PayloadPerSecond  *float64 `json:"payloadPerSecond"`

	MessagesHistory *[]HistoryEntry `json:"messagesHistory"`
}

type HistoryEntry struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type Producer struct {
	Component `json:",inline"`
	Protocol  string `json:"protocol"`

	MessagesPerSecondConfigured *float64 `json:"messagesPerSecondConfigured"`
	MessagesPerSecondScheduled  *float64 `json:"messagesPerSecondScheduled"`
	MessagesPerSecondSent       *float64 `json:"messagesPerSecondSent"`
	MessagesPerSecondFailed     *float64 `json:"messagesPerSecondFailed"`
	MessagesPerSecondErrored    *float64 `json:"messagesPerSecondErrored"`

	ConnectionsConfigured  *float64 `json:"connectionsConfigured"`
	ConnectionsEstablished *float64 `json:"connectionsEstablished"`

	ChartData   []ChartEntry       `json:"chartData"`
	ChartLegend []ChartLegendEntry `json:"chartLegend"`
}

type ChartEntry struct {
	Key   string  `json:"x"`
	Value float64 `json:"y"`
}
type ChartLegendEntry struct {
	Name string `json:"name"`
}
