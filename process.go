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
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/anthonynsimon/bild/transform"
	"github.com/korandiz/v4l"
	"github.com/labstack/gommon/color"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
	"golang.org/x/crypto/ssh/terminal"
)

func process(cam *v4l.Device) {
	// Initialize image counter
	i := 1

	// Initialize warning error
	j := 0

	// Infinite processing loop
	for {
		// Generate ID for logging
		ID := strconv.Itoa(i)

		// Start processing time
		start := time.Now()

		var imagePath string
		if arguments.Keep == true {
			imagePath = arguments.Output + "/" + start.Format("2006-01-02-15-04-05") + ".jpg"
		} else {
			imagePath = start.Format("2006-01-02-15-04-05")
		}

		// Read frame from camera
		buf, err := cam.Capture()
		if err != nil {
			log.Println("Capture:", err)
			proc, _ := os.FindProcess(os.Getpid())
			proc.Signal(os.Interrupt)
			break
		}

		// Decode frame to jpeg
		//buf.Seek(0, 0)
		img, err := jpeg.DecodeIntoRGBA(buf, &jpeg.DecoderOptions{})
		if err != nil {
			logError("Error decoding jpeg image: "+err.Error(), "[ERROR]")
			os.Exit(1)
		}

    // Keep img on disk
    if arguments.Keep == true {
      go keepImg(imagePath, img)
    }

		// Show warning if model size can't be used for capture
		im, _, err := image.DecodeConfig(buf)
		if (im.Width != arguments.Width || im.Height != arguments.Height) && j == 0 {
			logError("Model size specified as parameters can't be use for capture with this camera.", "[WARNING]")
			logError("Input image will be resized during DeepDetect processing.", "[WARNING]")
			j = 1
		}

		// Flip image
		if arguments.Mirror == true {
			img = transform.FlipH(img)
		}

		// Encode as base64
		buffer64 := new(bytes.Buffer)
		err = jpeg.Encode(buffer64, img, &jpeg.EncoderOptions{Quality: 100})
		if err != nil {
			log.Fatal("Error encoding image to base64: " + err.Error())
			os.Exit(1)
		}
		imageBase64 := base64.StdEncoding.EncodeToString(buffer64.Bytes())

		/*if arguments.Verbose == "DEBUG" {
			logSuccess(color.Cyan("Saving frame to ")+
				color.Yellow(imagePath), "["+
				ID+"] [DEBUG]")
		}

		if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
			logSuccess("Image "+imagePath+" captured", "["+
				ID+"] [INFO]")

			logSuccess("Sending to DeepDetect for processing", "["+
				ID+"] [INFO]")
		}*/

		// Process image with DeepDetect
		deepdetectProcess(imagePath, ID, img, start, imageBase64)

		if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
			// Log processing time
			logSuccess("Picture processed in "+
				color.Cyan(time.Since(start)), "["+
				ID+"] [INFO]")
		}

		// Increment image counter
		i++

		// Get terminal size and print pretty horizontal bar
		width, _, err := terminal.GetSize(1)
		if err != nil {
			logError("Unable to get terminal dimensions!", "[ERROR]")
		}

		// Show log message about waiting period
		if arguments.Waiting > 0 {
			logSuccess("Waiting " + strconv.Itoa(arguments.Waiting) + " seconds before next request", "[INFO]")
		}

		// Pretty horizontal bar displaying
		if arguments.Verbose == "INFO" || arguments.Verbose == "DEBUG" {
			for j := 0; j < width; j++ {
				fmt.Print(color.Green("="))
			}
		}

		// Wait `arguments.Waiting` seconds before next request
		if arguments.Waiting > 0 {
			time.Sleep(time.Duration(arguments.Waiting) * time.Second)
		}
	}
}
