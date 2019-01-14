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
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uilive"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/labstack/gommon/color"
	live "github.com/saljam/mjpeg"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Declare wrtier for one-line logging
	writer = uilive.New()

	// Declare stream for web preview
	stream = live.NewStream()

	// InfluxDB variables
	influxConnection = false
	influxClient     client.Client

	// Elapsed time logging
	elapsedTimes []time.Duration

	// HTTP settings
	httpTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
	}
	httpClient = &http.Client{
		Timeout:   time.Second * 60,
		Transport: httpTransport,
	}
)

func init() {
	/*
	   Safety net for 'too many open files' issue on legacy code.
	   Set a sane timeout duration for the http.DefaultClient, to ensure idle connections are terminated.
	   Reference: https://stackoverflow.com/questions/37454236/net-http-server-too-many-open-files-error
	*/
	http.DefaultClient.Timeout = time.Minute * 10
}

func main() {
	var cam *v4l.Device
	var err error

	// Parse arguments
	argumentsParsing(os.Args)

	// InfluxDB
	if arguments.Influx == true {
		initDB()
	}

	// If --detection is triggered, set --best to 1
	if arguments.Detection == true {
		arguments.Best = 1
	}

	// Service creation
	if arguments.SSL == true && arguments.Create == true {
		logSuccess("Creating service..", "[INFO]")
		createService("https://" + arguments.Host + ":" + arguments.Port)
	} else if arguments.Create == true {
		logSuccess("Creating service..", "[INFO]")
		createService("http://" + arguments.Host + ":" + arguments.Port)
	}

	// Start capture
	logSuccess("Starting capture", "[INFO]")
	devicePath := "/dev/video" + strconv.Itoa(arguments.DeviceID)
	logSuccess("Device: "+devicePath, "[INFO]")

	// Open camera
	cam, err = v4l.Open(devicePath)
	if err != nil {
		logError("Error opening video capture device", "[FATAL]")
		os.Exit(1)
	}

	// Set camera properties
	// Fetch config
	cfg, err := cam.GetConfig()
	if err != nil {
		logError(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Set parameters
	cfg.Format = mjpeg.FourCC
	cfg.Width = arguments.Width
	cfg.Width = arguments.Width

	// Apply config
	err = cam.SetConfig(cfg)
	if err != nil {
		logError(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Turn on cam
	err = cam.TurnOn()
	if err != nil {
		logError(err.Error(), "[ERROR]")
		os.Exit(1)
	}

	// Verify config
	cfg, err = cam.GetConfig()
	if err != nil {
		logError(err.Error(), "[ERROR]")
		os.Exit(1)
	}
	if cfg.Format != mjpeg.FourCC {
		logError("Failed to set MJPEG format.", "[ERROR]")
		os.Exit(1)
	}

	// Info message for web preview
	if arguments.Preview != "" {
		logSuccess("Starting web preview on "+arguments.Preview, "[INFO]")
		go http.Handle("/", stream)
		go http.ListenAndServe(arguments.Preview, nil)
	}

	// Get terminal size and print pretty horizontal bar
	width, _, err := terminal.GetSize(1)
	if err != nil {
		logError("Unable to get terminal dimensions!", "[ERROR]")
	}

	// Pretty horizontal bar displaying
	for j := 0; j < width; j++ {
		fmt.Print(color.Green("="))
	}

	// Start listening for updates and render
	if arguments.Verbose != "DEBUG" && arguments.Verbose != "INFO" {
		writer.Start()
	}

	// Start image processing
	process(cam)
}
