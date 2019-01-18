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
	"github.com/ctron/iot-simulator-console/pkg/metrics"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"os"

	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	flag.Parse()

	namespace, _ := os.LookupEnv("NAMESPACE")

	promcfg, err := metrics.BuildConfiguration(namespace)
	if err != nil {
		log.Fatalf("Unable to build metrics configuration: %v", err)
	}
	metricsClient, err := metrics.NewMetrics(promcfg)
	if err != nil {
		log.Fatalf("Unable to build metrics configuration: %v", err)
	}

	log.Printf("Using Prometheus endpoint: %s", promcfg.Url)

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

	router := gin.Default()

	controller := data.NewController(namespace, client, appsclient, metricsClient)

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
