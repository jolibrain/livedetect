/**
 * LiveDetect
 * Copyright (c) 2018 Jolibrain
 * Author: Corentin Barreau <corentin.barreau@epitech.eu>
 *
 * This file is part of livedetect.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in livedetect without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of livedetect, and to permit persons to whom livedetect is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of livedetect.
 *
 * LIVEDETECT IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH LIVEDETECT OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"os"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/gommon/color"
	"golang.org/x/crypto/ssh/terminal"
)

func queryDB(DB string, clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: DB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func writePoints(clnt client.Client, category string, probability string) {
	// Create batch point
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: arguments.InfluxDB,
	})
	if err != nil {
		logError("Error creating batch points: "+err.Error(), "[ERROR]")
	}

	// Declare tags
	tags := map[string]string{
		"category": category,
	}

	fields := map[string]interface{}{
		"probability": probability,
	}
	// Create point
	pt, err := client.NewPoint(
		"objects",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		logError("Error creating data point: "+err.Error(), "[ERROR]")
	}
	bp.AddPoint(pt)

	// Write points InfluxDB
	if err := clnt.Write(bp); err != nil {
		logError("Error pushing data to InfluxDB: "+err.Error(), "[ERROR]")
	}
}

func initDB() {
	// Get terminal size and print pretty horizontal bar
	width, _, err := terminal.GetSize(1)
	if err != nil {
		logError("Unable to get terminal dimensions!", "[ERROR]")
	}

	// Pretty horizontal bar displaying
	for j := 0; j < width; j++ {
		fmt.Print(color.Green("="))
	}

	logSuccess("Initializing "+color.Cyan("InfluxDB")+color.Yellow(" connection.."), "[INFO]")
	logSuccess("Host: "+color.Cyan(arguments.InfluxHost), "[INFO]")
	logSuccess("User: "+color.Cyan(arguments.InfluxUser), "[INFO]")
	logSuccess("Pass: "+color.Cyan(arguments.InfluxPass), "[INFO]")
	logSuccess("Database: "+color.Cyan(arguments.InfluxDB), "[INFO]")

	// Initialize HTTP connection to InfluxDB
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     arguments.InfluxHost,
		Username: arguments.InfluxUser,
		Password: arguments.InfluxPass,
		Timeout:  time.Second * 5,
	})
	if err != nil {
		logError("Error initializing "+color.Cyan("InfluxDB")+
			color.Yellow(" HTTP client: ")+
			color.Red(err.Error()), "[ERROR]")
		os.Exit(1)
	}
	defer client.Close()
	logSuccess("Connected to "+color.Cyan("InfluxDB")+
		color.Yellow(" instance."), "[INFO]")

	// Initialize database if it doesn't already exist
	_, err = queryDB(arguments.InfluxDB, client, fmt.Sprintf("CREATE DATABASE %s", arguments.InfluxDB))
	if err != nil {
		logError("Error while creating the database "+
			arguments.InfluxDB+": "+err.Error(), "[ERROR]")
		os.Exit(1)
	}
	logSuccess("Database "+color.Cyan(arguments.InfluxDB)+
		color.Yellow(" ready."), "[INFO]")

	// Pretty horizontal bar displaying
	for j := 0; j < width; j++ {
		fmt.Print(color.Green("="))
	}

	influxConnection = true
	influxClient = client
}
