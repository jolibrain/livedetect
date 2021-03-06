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
	"os"
	"strconv"
	"strings"
	"time"
	"bufio"
	"image"

  "encoding/json"
  "io/ioutil"

	"github.com/jolibrain/godd"
	jpeg "github.com/pixiv/go-libjpeg/jpeg"
)

func keepImg(filePath string, img image.Image) {

  f, err := os.Create(filePath)
  if err != nil {
    logError("Error creating file: " + err.Error(), "[ERROR]")
		panic(err.Error())
  }

  b := bufio.NewWriter(f)
  defer func() {
    b.Flush()
    f.Close()
  }()

  if err := jpeg.Encode(b, img, &jpeg.EncoderOptions{Quality: 100}); err != nil {
    panic(err)
  }
  return
}

func keepJson(filePath string, result godd.PredictResult) {
  file, _ := json.Marshal(result)
  _ = ioutil.WriteFile(filePath, file, 0644)
}

func floatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

// AverageTime calculate the average elapsed time
// from a slice of time.Duration
func averageTime(elapsedTimes []time.Duration) time.Duration {
	var totalTime time.Duration
	for _, value := range elapsedTimes {
		totalTime += value
	}

	return totalTime / time.Duration(len(elapsedTimes))
}

// Contains check if a string is contained inside a string slice
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
