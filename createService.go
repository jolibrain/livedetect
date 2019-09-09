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

	"github.com/jolibrain/godd"
)

func createService(URL string) {

  // Create service struct for service creation
  // parameters
  var service godd.ServiceRequest

  if arguments.ServiceConfig.Create == nil {

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

    if len(arguments.MlLibDataType) > 0 {
      service.Parameters.Mllib.Datatype = arguments.MlLibDataType
    }

    if arguments.MlLibMaxBatchSize != -1 {
      service.Parameters.Mllib.MaxBatchSize = arguments.MlLibMaxBatchSize
    }

    if arguments.MlLibMaxWorkspaceSize  != -1 {
      service.Parameters.Mllib.MaxWorkspaceSize = arguments.MlLibMaxWorkspaceSize
    }

    if service.Model.Init != "" {
       service.Model.CreateRepository = true
    }

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

    // Delete service if already existing
    godd.DeleteService(URL, service.Name)

    // Send the service creation request
    creationResult, err := godd.CreateService(URL, &service)
    if err != nil {
      log.Fatal(err)
    }

        // Error handling
    log.Println("creation status=",creationResult.Status.Code)
    if creationResult.Status.Code != 201 {
      logError("Unable to create service!", "[ERROR]")
      logError("Error: "+creationResult.Status.Msg, "[ERROR]")
      //os.Exit(1)
    } else {
      logSuccess("Service created!", "[INFO]")
    }

  } else {

    // Iterate through arguments.ServiceConfig predict services
    for i := 0; i < len(arguments.ServiceConfig.Create); i++ {

      service = arguments.ServiceConfig.Create[i]

      // Delete service if already existing
      godd.DeleteService(URL, service.Name)

      // Send the service creation request
      creationResult, err := godd.CreateService(URL, &service)
      if err != nil {
        log.Fatal(err)
      }

          // Error handling
      log.Println("creation status=",creationResult.Status.Code)
      if creationResult.Status.Code != 201 {
        logError("Unable to create service!", "[ERROR]")
        logError("Error: "+creationResult.Status.Msg, "[ERROR]")
        //os.Exit(1)
      } else {
        logSuccess("Service created!", "[INFO]")
      }

    }
  }
}
