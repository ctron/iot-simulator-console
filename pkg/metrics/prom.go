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

package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/ctron/iot-simulator-console/pkg/data"

	"github.com/prometheus/common/log"

	"os"

	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"
)

type MetricsClient struct {
	api promv1.API
}

type Configuration struct {
	Url string
}

func BuildConfiguration(namespace string) (Configuration, error) {

	prometheusUrl := os.Getenv("PROMETHEUS_URL")
	if prometheusUrl == "" {
		prometheusHost := os.Getenv("PROMETHEUS_HOST")
		if prometheusHost == "" && namespace != "" {
			prometheusHost = "prometheus-operated." + namespace + ".svc"
		}
		prometheusPort := os.Getenv("PROMETHEUS_PORT")
		if prometheusPort == "" {
			prometheusPort = ":9090"
		} else {
			prometheusPort = ":" + prometheusPort
		}
		prometheusProto := os.Getenv("PROMETHEUS_PROTO")
		if prometheusProto == "" {
			prometheusProto = "http"
		}
		if prometheusProto != "" && prometheusHost != "" {
			prometheusUrl = prometheusProto + "://" + prometheusHost + prometheusPort
		}
	}

	if prometheusUrl == "" {
		return Configuration{}, fmt.Errorf("unable to build configuration")
	}

	return Configuration{
		prometheusUrl,
	}, nil
}

func NewMetrics(configuration Configuration) (*MetricsClient, error) {

	promClient, err := promapi.NewClient(promapi.Config{Address: configuration.Url})
	if err != nil {
		return nil, err
	}

	promApi := promv1.NewAPI(promClient)

	result := MetricsClient{
		api: promApi,
	}

	return &result, nil

}

func (c *MetricsClient) Query(ctx context.Context, query string) (*prommodel.Value, error) {
	log.Info("Query: ", query)

	s := time.Now()

	val, err := c.api.Query(ctx, query, s)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (c *MetricsClient) QueryRange(ctx context.Context, query string, duration time.Duration, step time.Duration) (*prommodel.Value, error) {
	log.Info("QueryRange: ", query)

	e := time.Now()
	s := e.Add(-duration)

	r := promv1.Range{Start: s, End: e, Step: step}

	val, err := c.api.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (c *MetricsClient) QueryMap(ctx context.Context, query string) (map[string]float64, error) {

	val, err := c.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	switch v := (*val).(type) {
	case *prommodel.Scalar:
		return nil, fmt.Errorf("missing map structure: %v", v)
	case prommodel.Vector:
		log.Info("Query result - vector: ", v)
		var result = make(map[string]float64)
		for _, e := range v {
			code := string(e.Metric["code"])
			result[code] = float64(e.Value)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown result type: %v / %v", (*val).Type().String(), (*val).String())
	}
}

func (c *MetricsClient) QuerySingle(ctx context.Context, query string) (*float64, error) {

	val, err := c.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	switch v := (*val).(type) {
	case *prommodel.Scalar:
		f := float64(v.Value)
		log.Info("Query result - scalar: ", f)
		return &f, nil
	case prommodel.Vector:
		log.Info("Query result - vector: ", v)
		if len(v) > 0 {
			f := float64(v[0].Value)
			return &f, nil
		} else {
			return nil, fmt.Errorf("missing values in vector result")
		}
	default:
		return nil, fmt.Errorf("unknown result type: %v / %v", (*val).Type().String(), (*val).String())
	}

}

func (c *MetricsClient) QueryArray(ctx context.Context, duration time.Duration, step time.Duration, query string, overrideLabel []string) (*[]data.HistoryEntry, error) {

	val, err := c.QueryRange(ctx, query, duration, step)
	if err != nil {
		return nil, err
	}

	switch v := (*val).(type) {

	case prommodel.Matrix:
		log.Info("Query result - matrix: ", v)

		l := 0
		for _, m := range v {
			l += len(m.Values)
		}

		f := make([]data.HistoryEntry, l)
		i := 0
		for x, m := range v {

			var label = m.Metric.String()
			if overrideLabel != nil && len(overrideLabel) > x {
				label = overrideLabel[x]
			}

			for _, sp := range m.Values {
				f[i].Name = label
				f[i].X = float64(sp.Timestamp.Unix())
				f[i].Y = float64(sp.Value)
				i++
			}
		}

		return &f, nil

	default:
		return nil, fmt.Errorf("unknown result type: %v / %v", (*val).Type().String(), (*val).String())
	}

}
