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
	"fmt"
	"image"
	"image/draw"
	"strconv"
	"time"
	"strings"

	"github.com/jolibrain/godd"
	"github.com/labstack/gommon/color"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
)

func printResponse(request godd.PredictRequest, result godd.PredictResult, ID string, img image.Image, filePath string, startTime time.Time) {

	if result.Status.Code != 200 {
    logError("Unexpected response status code: " + strconv.Itoa(result.Status.Code), "[ERROR]")
    return
  }

	// Initialize RGBA image
	var imgRGBA *image.RGBA

	// Initialize categories counter
	categories := 0

	// Print complete classes array if --verbose is true
	if arguments.Verbose == "DEBUG" {
		logSuccess("Raw classes: "+color.Cyan(result.Body.Predictions[0]),
			"["+ID+"] [DEBUG]")
	}

	// Iterate through classes and print cat, probs and bbox if
	// --detection is used
	if len(result.Body.Predictions) != 0 {
		// Convert to RGBA
		b := img.Bounds()
		imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

		// Loop over predictions
		for _, class := range result.Body.Predictions[0].Classes {
			cat := class.Cat
			// Classes selection statement
			if arguments.SelectClasses == true && sliceContains(*arguments.Classes, cat) == false {
			} else {
				if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
					fmt.Print("\n")
					logSuccess("Category: "+color.Cyan(cat),
						"["+ID+"] [INFO]")
				}
				prob := class.Prob
				if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
					logSuccess("Probability: "+color.Cyan(prob),
						"["+ID+"] ------")
				}

				// Print bbox if --detection
				if arguments.Detection == true {
					bbox := class.Bbox
					if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
						logSuccess("Coordinates: "+color.Cyan(bbox),
							"["+ID+"] ------")
					}
					// Parse mask output
					if arguments.Mask == true {
						mask := class.Mask
						if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
							logSuccess("Mask Format: "+color.Cyan(mask.Format),
								"["+ID+"] ------")
							logSuccess("Mask Width: "+color.Cyan(mask.Width),
								"["+ID+"] ------")
							logSuccess("Mask Height: "+color.Cyan(mask.Height),
								"["+ID+"] ------")
							fmt.Print("\n")
						}
						if arguments.Preview != "" || arguments.Keep == true {
							imgRGBA = writeMask(img, result, categories, ID)
							img = imgRGBA
						}
					} else {
						if arguments.Preview != "" || arguments.Keep == true {
							imgRGBA = writeBoundingBox(img, result, categories, ID)
							img = imgRGBA
						}
					}
				}

				// Push data to InfluxDB
				if influxConnection == true {
					if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
						logSuccess("Sending gathered data to "+
							color.Cyan("InfluxDB"), "[INFO]")
					}
					go writePoints(influxClient, cat, floatToString(prob))
					if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
						logSuccess("Data succesfully pushed to "+
							color.Cyan("InfluxDB"), "[INFO]")
					}
				}
				categories++
			}
		}
	}

	// Update real time buffer
	if arguments.Verbose != "DEBUG" && arguments.Verbose != "INFO" {
		elapsed := time.Since(startTime)
		elapsedTimes = append(elapsedTimes, elapsed)
		fmt.Fprintln(writer, color.Green("[✔] ["+ID+"]")+
			color.Yellow(" Image processed! ")+
			color.Green(strconv.Itoa(categories))+
			color.Yellow(" categories found! ")+
			color.Cyan(elapsed)+
			color.Green(" | ")+
			color.Yellow("Average processing time: ")+
			color.Cyan(averageTime(elapsedTimes)))
	}

	// Preview window
	if arguments.Preview != "" {
		// Convert image to buffer
		buf := new(bytes.Buffer)
		if imgRGBA != nil {
			err := jpeg.Encode(buf, imgRGBA, &jpeg.EncoderOptions{Quality: 50})
			if err == nil {
				go stream.UpdateJPEG(buf.Bytes())
			} else {
				logError("Can't encode frame to live stream.", "[ERROR]")
			}
		}
	}

  // Keep json on disk
  if arguments.Keep == true {

    // Place json file next to processed image file
    var logPath string
    logPath = strings.TrimSuffix(filePath, ".jpg") + ".json"

    // Add Service name suffix if specified
    if request.Service != "" {
      logPath = strings.Replace(logPath, ".json", "_" + request.Service + ".json", -1)
    }

    // Write predict response inside json file
    go keepJson(logPath, result)

  }
}
