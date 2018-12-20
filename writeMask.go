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
	"math/rand"

	"github.com/CorentinB/gobbox"
	"github.com/CorentinB/godd"
)

func writeSquare(red, green, blue uint8, x, y int, img draw.Image) {
	// Set color
	c := color.RGBA{
		red, green, blue, 1,
	}

	// Write square
	img.Set(x, y, c)
	img.Set(x+1, y+1, c)
	img.Set(x-1, y-1, c)
	img.Set(x+1, y-1, c)
	img.Set(x-1, y+1, c)

	img.Set(x+1, y+2, c)
	img.Set(x-1, y-2, c)
	img.Set(x+1, y-2, c)
	img.Set(x-1, y+2, c)

	img.Set(x+2, y+2, c)
	img.Set(x-2, y-2, c)
	img.Set(x+2, y-2, c)
	img.Set(x-2, y+2, c)

	img.Set(x+2, y+1, c)
	img.Set(x-2, y-1, c)
	img.Set(x+2, y-1, c)
	img.Set(x-2, y+1, c)

	img.Set(x, y+1, c)
	img.Set(x, y-1, c)
	img.Set(x, y-1, c)
	img.Set(x, y+1, c)
	img.Set(x, y+2, c)
	img.Set(x, y-2, c)
	img.Set(x, y-2, c)
	img.Set(x, y+2, c)

	img.Set(x+1, y, c)
	img.Set(x-1, y, c)
	img.Set(x+1, y, c)
	img.Set(x-1, y, c)
	img.Set(x+2, y, c)
	img.Set(x-2, y, c)
	img.Set(x+2, y, c)
	img.Set(x-2, y, c)
}

func writeMask(img image.Image, result godd.PredictResult, class int, ID string) (imgRGBA *image.RGBA) {
	// Convert to RGBA
	b := img.Bounds()
	imgRGBA = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), img, b.Min, draw.Src)

	if arguments.SelectClasses == true && sliceContains(*arguments.Classes, result.Body.Predictions[0].Classes[class].Cat) == false {
	} else {
		var size image.Point
		size.X = arguments.Width
		size.Y = arguments.Height
		maskDataRaw := result.Body.Predictions[0].Classes[class].Mask.Data

		i := 0
		// Generate a random color for filling mask
		red := uint8(rand.Intn(255))
		green := uint8(rand.Intn(255))
		blue := uint8(rand.Intn(255))

		// Positions for writing mask
		xStart := int(result.Body.Predictions[0].Classes[class].Bbox.Xmin)
		length := int(result.Body.Predictions[0].Classes[class].Bbox.Xmin) +
			result.Body.Predictions[0].Classes[class].Mask.Width
		yStart := int(result.Body.Predictions[0].Classes[class].Bbox.Ymin)
		height := int(result.Body.Predictions[0].Classes[class].Bbox.Ymin) +
			result.Body.Predictions[0].Classes[class].Mask.Height

		// Loop though all the x
		for y := yStart; y < height; y++ {
			// And now loop through all of this x's y
			for x := xStart; x < length; x++ {
				if arguments.Contour == true {
					if i != 0 && i != len(maskDataRaw)-1 {
						if maskDataRaw[i] == 1 && maskDataRaw[i-1] == 0 ||
							maskDataRaw[i] == 1 && maskDataRaw[i+1] == 0 {
							writeSquare(red, green, blue, x, y, imgRGBA)
						}
					}
					i++
				} else {
					if maskDataRaw[i] == 1 {
						c := color.RGBA{
							red, green, blue, 1,
						}
						imgRGBA.Set(x, y, c)
					}
					i++
				}
			}
		}

		// Write category
		bbox := result.Body.Predictions[0].Classes[class].Bbox
		// Get box corners
		left := int(bbox.Xmin)
		top := int(bbox.Ymax)
		gobbox.WriteLabel(imgRGBA, result.Body.Predictions[0].Classes[class].Cat, left, top, color.RGBA{255, 255, 255, 255})
	}
	return imgRGBA
}
