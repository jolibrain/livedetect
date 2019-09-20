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
	"bytes"
	"image"
  "time"

	"github.com/jolibrain/godd"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
)

func deepdetectProcess(imagePath string, ID string, img image.Image, startTime time.Time, imageBase64 string) {

	var predictURL string
  var request godd.PredictRequest

	if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
		logSuccess("Processing image "+imagePath,
			"["+ID+"] [INFO]")
	}

	// Generate predict URL
	if arguments.SSL == true {
		predictURL = "https://" + arguments.Host + ":" + arguments.Port
	} else {
		predictURL = "http://" + arguments.Host + ":" + arguments.Port
	}

  if arguments.Path != "" {
    predictURL = predictURL + arguments.Path
  }

	// Execute predict
  if arguments.ServiceConfig.Predict == nil {

    // Use only arguments.Service as predict service
    response := predict(predictURL, imageBase64, ID)

    // Handle response
    printResponse(request, response, ID, img, imagePath, startTime, 0)

  } else {

    // Iterate through arguments.ServiceConfig predict services
    for i := 0; i < len(arguments.ServiceConfig.Predict); i++ {

      request = arguments.ServiceConfig.Predict[i]

      if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
        logSuccess("Request on service " + request.Service,
          "["+ID+"] [INFO]")
      }

      response := predictWithRequest(request, predictURL, imageBase64, ID)

      // Handle response
      img = printResponse(request, response, ID, img, imagePath, startTime, i)

    }

    // Preview window
    if arguments.Preview != "" {
      // Convert image to buffer
      buf := new(bytes.Buffer)
      if img != nil {
        err := jpeg.Encode(buf, img, &jpeg.EncoderOptions{Quality: 50})
        if err == nil {
          go stream.Update(buf.Bytes())
        } else {
          logError("Can't encode frame to live stream.", "[ERROR]")
        }
      }
    }


  }

}
