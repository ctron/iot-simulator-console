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

package main

import (
	"flag"
	"github.com/ctron/iot-simulator-console/pkg/data"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"

	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	promapi "github.com/prometheus/client_golang/api"
)

func main() {
	flag.Parse()

	namespace, _ := os.LookupEnv("NAMESPACE")

	prometheusUrl := os.Getenv("PROMETHEUS_URL")
	if prometheusUrl == "" {
		prometheusHost := os.Getenv("PROMETHEUS_HOST")
		if prometheusHost == "" {
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
			prometheusPort = "http"
		}
		prometheusUrl = prometheusProto + "://" + prometheusHost + prometheusPort
	}
	log.Printf("Using Prometheus endpoint: %s", prometheusUrl)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("Failed to get configuration: %s", err)
		os.Exit(1)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes client: %v", err.Error())
	}

	appsclient, err := appsv1.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes client: %v", err.Error())
	}

	promClient, err := promapi.NewClient(promapi.Config{Address: prometheusUrl})
	promApi := v1.NewAPI(promClient)

	controller := data.NewController(namespace, client, appsclient, promApi)
	router := gin.Default()

	router.Use(
		static.Serve(
			"/",
			static.LocalFile("./build", true),
		),
	)

	// Setup route group for the API
	api := router.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// api.GET("/overview", JokeHandler)
	api.GET("/overview", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		result, err := controller.BuildOverview()
		if err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, result)
		}
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error running router: %v", err)
	}
}

func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Jokes handler not implemented yet",
	})
}
