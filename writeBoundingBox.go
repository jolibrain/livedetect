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

func writeBoundingBox(img image.Image, result godd.PredictResult, class int, ID string, index int) (imgRGBA *image.RGBA) {

	// Convert to RGBA
	b := img.Bounds()
	imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

  // Set colors for bbox
  boxColor := color.RGBA{228,26,28,255}
  switch index {
  case 1:
    boxColor = color.RGBA{55,126,184,255}
  case 2:
    boxColor = color.RGBA{77,175,74,255}
  case 3:
    boxColor = color.RGBA{152,78,163,255}
  case 4:
    boxColor = color.RGBA{255,127,0,255}
  case 5:
    boxColor = color.RGBA{255,255,51,255}
  case 6:
    boxColor = color.RGBA{166,86,40,255}
  case 7:
    boxColor = color.RGBA{247,129,191,255}
  case 8:
    boxColor = color.RGBA{153,153,153,255}
  }

  // Set color for label
	white := color.RGBA{255, 255, 255, 255}

	// Set coordinates for bbox
	x1 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmin)
	x2 := int(result.Body.Predictions[0].Classes[class].Bbox.Xmax)
	y1 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymin)
	y2 := int(result.Body.Predictions[0].Classes[class].Bbox.Ymax)

	// Draw the bounding box
	gobbox.DrawBoundingBox(imgRGBA, result.Body.Predictions[0].Classes[class].Cat, x1, x2, y2, y1, boxColor, white)

	return imgRGBA
}
