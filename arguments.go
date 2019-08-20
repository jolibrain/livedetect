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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"io/ioutil"
  "encoding/json"

	"github.com/akamensky/argparse"
	"github.com/jolibrain/godd"
)


var arguments = struct {
	// InfluxDB flags
	Influx     bool
	InfluxHost string
	InfluxUser string
	InfluxPass string
	InfluxDB   string
	// OpenCV flags
	FPS      float64
	DeviceID int
	Video    string
	Mirror   bool
	// Livedetect flags
	Output        string
	Keep          bool
	Preview       string
	Picamera      bool
	Verbose       string
	SelectClasses bool
	Classes       *[]string
	// Predict flags
	Host       string
	Port       string
	Path       string
	Width      int
	Height     int
	Best       int
	Service    string
	ServiceConfig []godd.PredictRequest
	Detection  bool
	Confidence float64
	SSL        bool
	Waiting	   int
	// Creation flags
	Create             bool
	GPU                bool
	Nclasses           int
	Template           string
	ModelRepository    string
	ServiceDescription string
	Mllib              string
	Connector          string
	Init               string
	MlLibDataType           string
	MlLibMaxBatchSize       int
	MlLibMaxWorkspaceSize   int
	// Mask
	Mask       bool
	Contour    bool
	Extensions *[]string
	Mean       *[]float64
}{
	// Default arguments
	Width:      227,
	Height:     227,
	Confidence: 0.10,
	FPS:        30.00,
	DeviceID:   0,
	Best:       3,
	Waiting:    0}

func argumentsParsing(args []string) {
	// Create new parser object
	parser := argparse.NewParser("LiveDetect", "Real-time image processing with DeepDetect")

	// Create flags
	contour := parser.Flag("", "contour", &argparse.Options{
		Required: false,
		Help:     "Draw contour of mask instead of full mask",
		Default:  false})

	host := parser.String("", "host", &argparse.Options{
		Required: false,
		Help:     "Host of your DeepDetect instance (i.e: localhost)",
		Default:  "localhost"})

	port := parser.String("p", "port", &argparse.Options{
		Required: false,
		Help:     "Port used by your DeepDetect instance",
		Default:  "8080"})

	path := parser.String("", "path", &argparse.Options{
		Required: false,
		Help:     "Url Path of your DeepDetect instance (i.e: /api/deepdetect)",
		Default:  ""})

	init := parser.String("", "init", &argparse.Options{
		Required: false,
		Help:     "Path of pre-made compressed archive for DeepDetect's model loading",
		Default:  ""})

	influx := parser.Flag("", "influx", &argparse.Options{
		Required: false,
		Help:     "If data should be logged to and InfluxDB databse",
		Default:  false})

	influxHost := parser.String("", "influx-host", &argparse.Options{
		Required: false,
		Help:     "Host used by your InfluxDB instance",
		Default:  "http://localhost:8086"})

	influxUser := parser.String("", "influx-user", &argparse.Options{
		Required: false,
		Help:     "User used to connect to your InfluxDB database",
		Default:  "admin"})

	influxPass := parser.String("", "influx-pass", &argparse.Options{
		Required: false,
		Help:     "Password used to connect to your InfluxDB database",
		Default:  "admin"})

	influxDB := parser.String("", "influx-db", &argparse.Options{
		Required: false,
		Help:     "Name of the database that should be created or used by LiveDetect",
		Default:  "livedetect"})

	output := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Your output folder for saving images captured",
		Default:  "./"})

	width := parser.Int("", "width", &argparse.Options{
		Required: false,
		Help:     "Width of images captured",
		Default:  227})

	height := parser.Int("", "height", &argparse.Options{
		Required: false,
		Help:     "Height of images captured",
		Default:  227})

	FPS := parser.Float("", "fps", &argparse.Options{
		Required: false,
		Help:     "FPS of the camera",
		Default:  30.00})

	video := parser.String("", "video", &argparse.Options{
		Required: false,
		Help:     "Path of the input video",
		Default:  ""})

	verbose := parser.Selector("v", "verbose", []string{"INFO", "DEBUG"}, &argparse.Options{
		Required: false,
		Help:     "Verbosity, INFO or DEBUG",
		Default:  "INFO"})

	detection := parser.Flag("d", "detection", &argparse.Options{
		Required: false,
		Help:     "Run a detection model and perform bounding boxes writing on captured frames",
		Default:  false})

	// Mask support
	mask := parser.Flag("", "mask", &argparse.Options{
		Required: false,
		Help:     "Process mask on output",
		Default:  false})

	extensions := parser.List("e", "extensions", &argparse.Options{
		Required: false,
		Help:     "Extensions to use for mask"})

	mean := parser.List("m", "mean", &argparse.Options{
		Required: false,
		Help:     "Mean values for mask usage"})

	mirror := parser.Flag("", "mirror", &argparse.Options{
		Required: false,
		Help:     "Flip image before preview and processing",
		Default:  false})

	confidence := parser.Float("", "confidence", &argparse.Options{
		Required: false,
		Help:     "Only returns classifications or detections with probability strictly above threshold",
		Default:  0.10})

	picamera := parser.Flag("", "picamera", &argparse.Options{
		Required: false,
		Help:     "Use picamera instead of OpenCV",
		Default:  false})

	keep := parser.Flag("k", "keep", &argparse.Options{
		Required: false,
		Help:     "If pictures should be deleted after processing or not",
		Default:  false})

	best := parser.Int("b", "best", &argparse.Options{
		Required: false,
		Help:     "Number of different classes DeepDetect should return",
		Default:  3})

	service := parser.String("s", "service", &argparse.Options{
		Required: false,
		Help:     "Name of the service that should be used for the prediction",
		Default:  "imageserv"})

	serviceConfigPath := parser.String("C", "service-config", &argparse.Options{
		Required: false,
		Help:     "Configuration file for service(s) predict request",
		Default:  ""})

	preview := parser.String("P", "preview", &argparse.Options{
		Required: false,
		Help:     "Serve live processed images stream on a specific adress",
		Default:  ""})

	deviceID := parser.Int("", "device-id", &argparse.Options{
		Required: false,
		Help:     "Camera device ID OpenCV should connect to for capturing images",
		Default:  0})

	SSL := parser.Flag("", "ssl", &argparse.Options{
		Required: false,
		Help:     "Use HTTPS instead of HTTP",
		Default:  false})

	waiting := parser.Int("", "waiting", &argparse.Options{
		Required: false,
		Help:     "Waiting X seconds between predict requests",
		Default:  0})

	// Services creation flags
	create := parser.Flag("", "create", &argparse.Options{
		Required: false,
		Help:     "Create service before starting real time prediction",
		Default:  false})

	nclasses := parser.Int("", "nclasses", &argparse.Options{
		Required: false,
		Help:     "Number of classes for service creation",
		Default:  -1})

	template := parser.String("", "template", &argparse.Options{
		Required: false,
		Help:     "Template absolute path for service creation",
		Default:  "."})

	modelRepository := parser.String("", "repository", &argparse.Options{
		Required: false,
		Help:     "Model's repository absolute path for service creation",
		Default:  ""})

	serviceDescription := parser.String("", "description", &argparse.Options{
		Required: false,
		Help:     "Description for service creation",
		Default:  "LiveDetect service"})

	mllib := parser.String("", "mllib", &argparse.Options{
		Required: false,
		Help:     "mllib for service creation",
		Default:  "caffe"})

	connector := parser.String("", "connector", &argparse.Options{
		Required: false,
		Help:     "Connector for service creation input",
		Default:  "image"})

	GPU := parser.Flag("", "gpu", &argparse.Options{
		Required: false,
		Help:     "If the GPU should be used or not",
		Default:  false})

	MlLibDataType := parser.String("", "mllib-datatype", &argparse.Options{
		Required: false,
		Help:     "Mllib data type during service creation (fp32, fp16)",
		Default:  ""})

	MlLibMaxBatchSize := parser.Int("", "mllib-max-batch-size", &argparse.Options{
		Required: false,
		Help:     "Mllib max batch size",
		Default:  -1})

	MlLibMaxWorkspaceSize := parser.Int("", "mllib-max-workspace-size", &argparse.Options{
		Required: false,
		Help:     "Mllib max workspace size, in Mo",
		Default:  -1})

	// Classes
	selectClasses := parser.Flag("", "select-classes", &argparse.Options{
		Required: false,
		Help:     "Trigger classes selection. Use -c arguments for specifying classes to show",
		Default:  false})

	classes := parser.List("c", "classes", &argparse.Options{
		Required: false,
		Help:     "Show only these classes, need to use --select-classes to trigger classes selection"})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Check if template, modelTemplates and modelRepository
	// are absolute path, if yes, log error then exit
	// also check if at least model repository flag is filled
	/*if *create == true && strings.HasPrefix(*modelRepository == "" && == true && *init == "" {
		logError("You need to specify a path for the model repository!", "[ERROR]")
		os.Exit(1)
	}
	if *create == true &&
		strings.HasPrefix(*modelRepository, "/") == false {
		logError("You have to specify an absolute path for models and templates!", "[ERROR]")
		os.Exit(1)
	}*/

	// Check if --creation is triggered, if related flags are
	// also filled
	if *create == true && *nclasses <= 0 {
		logError("You have to specify a number of classes for service creation!", "[ERROR]")
		os.Exit(1)
	}

	// Convert path parameters to absolute paths
	var outputFolder, videoPath string
	if *output != "" {
		outputFolder, _ = filepath.Abs(*output)
	}
	if *video != "" {
		videoPath, _ = filepath.Abs(*video)
	}

	// Finally save the collected flags
	arguments.Host = *host
	arguments.Port = *port
	arguments.Path = *path
	arguments.Init = *init
	arguments.Width = *width
	arguments.Height = *height
	arguments.Mirror = *mirror
	arguments.FPS = *FPS
	arguments.Output = outputFolder
	arguments.Keep = *keep
	arguments.Service = *service

  if *serviceConfigPath != "" {

    // Open serviceConfigPath jsonFile
    jsonFile, err := os.Open(*serviceConfigPath)

    // if we os.Open returns an error then handle it
    if err != nil {
      logError("Can't open serviceConfig json",
        "[ERROR]")
      logError("Reason: "+err.Error(),
        "[ERROR]")
    }

    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, &arguments.ServiceConfig)

  }

	arguments.Best = *best
	arguments.Picamera = *picamera
	arguments.Preview = *preview
	arguments.Detection = *detection
	arguments.Confidence = *confidence
	arguments.Verbose = *verbose
	arguments.DeviceID = *deviceID
	arguments.SSL = *SSL
	arguments.Waiting = *waiting
	arguments.Create = *create
	arguments.Nclasses = *nclasses
	arguments.Template = *template
	arguments.ModelRepository = *modelRepository
	arguments.ServiceDescription = *serviceDescription
	arguments.Mllib = *mllib
	arguments.Connector = *connector
	arguments.GPU = *GPU
	arguments.MlLibDataType = *MlLibDataType
	arguments.MlLibMaxBatchSize = *MlLibMaxBatchSize
	arguments.MlLibMaxWorkspaceSize = *MlLibMaxWorkspaceSize
	arguments.Video = videoPath
	arguments.SelectClasses = *selectClasses
	arguments.Classes = classes
	// Mask support
	arguments.Mask = *mask
	arguments.Contour = *contour
	arguments.Extensions = extensions

	// Turn string slice into float64 slice
	if mean != nil {
		var meanFloat []float64
		for _, arg := range *mean {
			if n, err := strconv.ParseFloat(arg, 64); err == nil {
				meanFloat = append(meanFloat, n)
			}
		}
		arguments.Mean = &meanFloat
	}

	// InfluxDB
	arguments.Influx = *influx
	arguments.InfluxHost = *influxHost
	arguments.InfluxUser = *influxUser
	arguments.InfluxPass = *influxPass
	arguments.InfluxDB = *influxDB

	if arguments.Mask == true {
	   arguments.Detection = true
	}
}
