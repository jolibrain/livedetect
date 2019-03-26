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
	"github.com/jolibrain/godd"
)

func predict(URL string, image string, ID string) godd.PredictResult {
	// Create predict structure
	var predict godd.PredictRequest

	// Specify values for your prediction
	predict.Service = arguments.Service
	predict.Parameters.Input.Width = arguments.Width
	predict.Parameters.Input.Height = arguments.Height
	predict.Parameters.Output.Best = arguments.Best
	predict.Parameters.Output.Bbox = arguments.Detection
	predict.Parameters.Output.ConfidenceThreshold = arguments.Confidence
	predict.Parameters.Mllib.GPU = arguments.GPU
	predict.Data = append(predict.Data, image)

	// Mask support
	if arguments.Mask == true {
		predict.Parameters.Output.Mask = true
	}

	predictResult, err := godd.Predict(URL, &predict)
	if err != nil {
		logError("Can't execute request!",
			"["+ID+"] [ERROR]")
		logError("Reason: "+err.Error(),
			"["+ID+"] [ERROR]")
	}

	return predictResult
}
