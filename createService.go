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
	"log"
	"os"

	"github.com/CorentinB/godd"
)

func createService(URL string) {
	// Create service struct for service creation
	// parameters
	var service godd.ServiceRequest

	// Fill the service structure
	service.Name = arguments.Service
	service.Mllib = arguments.Mllib
	service.Parameters.Input.Connector = arguments.Connector
	service.Parameters.Input.Width = arguments.Width
	service.Parameters.Input.Height = arguments.Height
	service.Model.Repository = arguments.ModelRepository
	service.Model.Init = arguments.Init
	service.Parameters.Mllib.Nclasses = arguments.Nclasses
	service.Parameters.Mllib.GPU = arguments.GPU

	// Mask support
	if arguments.Mask == true {
		service.Mllib = "caffe2"
	}

	// Mask support
	if arguments.Mask == true {
		if len(*arguments.Extensions) == 0 {
			logError("You need to specify at least one extension for mask with the -e argument.", "[ERROR]")
			os.Exit(1)
		}
		arguments.Connector = "image"
		arguments.Mllib = "caffe2"
		service.Parameters.Input.Mean = *arguments.Mean
		service.Model.Extensions = *arguments.Extensions
	}

	// Send the service creation request
	creationResult, err := godd.CreateService(URL, &service)
	if err != nil {
		log.Fatal(err)
	}

	// Error handling
	if creationResult.Status.Code == 500 {
		logError("Unable to create service!", "[ERROR]")
		logError("Error: "+creationResult.Status.Msg, "[ERROR]")
	}

	// Log success
	if creationResult.Status.Code == 200 {
		logSuccess("Service created!", "[INFO]")
	}
}
