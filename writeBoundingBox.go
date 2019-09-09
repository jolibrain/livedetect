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
	"image"
	"image/color"
	"image/draw"

	"github.com/CorentinB/gobbox"
	"github.com/jolibrain/godd"
)

func writeBoundingBox(img image.Image, result godd.PredictResult, class int, ID string) (imgRGBA *image.RGBA) {

	// Convert to RGBA
	b := img.Bounds()
	imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

  // do not modify image if flag --keep-raw is true
  if arguments.KeepRaw == true {
    return imgRGBA
  }

	// Set colors for bbox
	red := color.RGBA{255, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Set coordinates for bbox
	x1 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmin)
	x2 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmax)
	y1 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymin)
	y2 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymax)

	// Draw the bounding box
	gobbox.DrawBoundingBox(imgRGBA, result.Body.Predictions[0].Classes[class].Cat, x1, x2, y2, y1, red, white)

	return imgRGBA
}
