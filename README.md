# LiveDetect

LiveDetect is a tool designed to easily process local video streams with Deep Learning models.

The code reads live imagery from a camera and processes every frame using [DeepDetect](https://github.com/jolibrain/deepdetect/).

LiveDetect and DeepDetect run on Desktop, any CPU, Nvidia's GPU, and Raspberry Pi3 (and other ARM devices) alike. 

Pre-trained Deep Learning [models](https://www.deepdetect.com/models/?opts={%22media%22:%22image%22,%22type%22:%22type-all%22,%22backend%22:[%22caffe%22,%22ncnn%22],%22platform%22:%22desktop%22,%22searchTerm%22:%22%22}#) are made available for [desktop](https://www.deepdetect.com/models/?opts={%22media%22:%22image%22,%22type%22:%22type-all%22,%22backend%22:[%22caffe%22,%22ncnn%22,%22caffe2%22,%22tensorflow%22],%22platform%22:%22desktop%22,%22searchTerm%22:%22%22}#) and [embedded systems like the Raspberry Pi](https://www.deepdetect.com/models/?opts={%22media%22:%22image%22,%22type%22:%22type-all%22,%22backend%22:[%22caffe%22,%22ncnn%22],%22platform%22:%22embedded%22,%22searchTerm%22:%22%22}#). 


Real-world use cases from DeepDetect customers with LiveDetect:

- Construction site safety and work monitoring

- Cars license plate OCR in parking lots

- Defect detection in manufactured precision pieces

![Example Traffic](example-traffic.gif)

## Set Up for Raspberry Pi 3

If you are using a **Raspberry Pi**, we made a step-by-step specific tutorial for the set up part [here](https://github.com/jolibrain/livedetect/wiki/Step-by-step-for-Raspberry-Pi-3). This tutorial is going to help you get your Raspberry Pi 3 ready. Once it's ready you can refer to the [examples section](https://github.com/jolibrain/livedetect#examples).

## Set Up for Desktop CPU and Nvidia's GPU

In order to use **LiveDetect**, you need to have a DeepDetect instance running and then you need to build LiveDetect. 

### DeepDetect instance 

You need a **DeepDetect** instance running and accessible from the machine where you want to use LiveDetect.

Install the dependencies for LiveDetect:

- `sudo apt-get install libjpeg-dev`

- create a `models` directory in your $HOME for the examples below to run without change.

- if you have a **Jetson ARM boards**

For Nvidia's Jetsons, you cannot use Docker and you should **build from source** as explained in the [quick start page](https://www.deepdetect.com/quickstart-server/?opts={%22os%22:%22ubuntu%22,%22source%22:%22build_source%22,%22compute%22:%22gpu%22,%22gpu%22:%22tx1%22,%22backend%22:[%22caffe%22],%22deepdetect%22:%22server%22}), using Caffe as a backend.

Once you are done with the quick start, you can directly [build LiveDetect from source](https://github.com/jolibrain/livedetect#buildfromsource).

- otherwise (no Jetson ARM board), **install Docker**

We are going to need Docker, as it works very well with very little overhead. If [you already have Docker installed](https://www.digitalocean.com/community/questions/how-to-check-for-docker-installation), you can jump to the next bullet point.

To install Docker:

```
curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
sudo groupadd docker
sudo usermod -aG docker $USER
```

- **Start a DeepDetect container**

**If you are on your CPU-only desktop**, you must run a CPU-only DeepDetect docker image:

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_cpu`

**If you are using an Nvidia GPU** (not for Jetson ARM boards), you need to install [nvidia-docker](https://github.com/NVIDIA/nvidia-docker). Once it's installed, you must run a GPU-enabled DeepDetect docker image:

- `nvidia-docker run -d -p 8080:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_gpu`

### Using LiveDetect

First check the [releases page](https://github.com/jolibrain/livedetect/releases):

- If the binary for your system is available:

You should download the appropriate version of LiveDetect and then make it executable using:

`chmod +x livedetect-$VERSION` replacing `$VERSION` by the version of LiveDetect you downloaded.

It's ready! You can directly jump to the [examples section](https://github.com/jolibrain/livedetect#examples).

- If you couldn't find the LiveDetect version fit for your device, you should build from source by following the next steps.

#### Requirements

In order to build from source and use LiveDetect, you need Go.

To install Go, download the approriate version for your OS [here](https://golang.org/dl/).

When the download is complete, go to the directory the file has been downloaded to and extract it to install the go toolchain:

- `tar -C $/usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz`

- Then make a local `go` directory in your `$HOME` with `mkdir $HOME/go`

- Set the goPATH with `export GOPATH=$HOME/go` and set path to the go binary with `export PATH=$PATH:$HOME/go/bin`. For more, refer to [this page](https://github.com/golang/go/wiki/SettingGOPATH).

#### LiveDetect build

Now that you have the required packages, go to your `$HOME` and fetch LiveDetect:

`go get -u github.com/jolibrain/livedetect`

Move in this directory, into to the LiveDetect file:

`cd ~/go/src/github.com/jolibrain/livedetect/`

From this directory, fetch all dependencies remaining:

`go get -u ./...`

Finally, build it!

`go build .`

You're ready for the following examples!

## Use-Cases

### Face detection + bounding boxes on web preview

This example starts LiveDetect and tells the DeepDetect instance listening from localhost:8080 to create a service named `voc`.

This service takes 300x300 frames and detect faces using the model fetched from the remote location specified by `--init`.

The preview is displayed on port 8888, that can be modified with flag `-P`.

V4L uses `--device-id` as the capture device (here device 0) and verbosity at the INFO level is triggered by `-v INFO`.

#### For Raspberry Pi

```
./livedetect-rpi3 \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/voc/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_voc_ncnn_300x300.tar.gz" \
    --confidence 0.3 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service voc \
    --nclasses 21
```

#### For Desktop

```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/voc/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_voc_ncnn_300x300.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service voc \
    --nclasses 21
```


#### For Nvidia GPU

```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib caffe \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/voc/ \
    --init "https://deepdetect.com/models/init/desktop/images/detection/detection_600.tar.gz" \
    --confidence 0.3 \
    -v INFO \
    --gpu \
    -P "0.0.0.0:8888" \
    --service voc \
    --nclasses 601
```


### Detection + mask

This command starts LiveDetect and tells the DeepDetect instance located at `localhost:8080` to create a service named `mask` for detection and mask preview.

This is for GPU machines only and uses Caffe2 [Detectron](https://github.com/facebookresearch/Detectron).

It uses files under the `--repository` path with extensions files under the path specified by `-e`.

Mean values are specified with `-m`, 81 classes, using the GPU, service should be configured for images of size 416x608 and predictions should be considered only if their confidence is at least 0.1.

Finally, using `--select-classes`, we specify classes we want to be previewed and processed, using classes names with the `-c` flags.

```
./livedetect \
    --port 8080 --host localhost \
    --width 416 --height 608 \
    --detection \
    --create --repository /home/user/test_mask \
    --init "https://www.deepdetect.com/models/init/desktop/images/detection/detectron_mask.tar.gz" \
    --mask -e /home/user/test_mask/mask \
    -m "102.9801" -m "115.9465" -m "122.7717" \
    --confidence 0.1 --nclasses 81 \
    --service mask \
    -P "0.0.0.0:8888" \
    --select-classes -c car -c person -c truck -c bike -c van
```

## Detection + InfluxDB

LiveDetect supports InfluxDB for storing and then displaying the results live with Grafana for instance. 

The categories detected and the probability for each category can be pushed to an InfluxDB database, here is the previous example, with the InfluxDB parameters.

To use InfluxDB, you first need to pass the `--influx` argument, then you can specify the host with `--influx-host`, the credentials that should be used with `--influx-user` and `--influx-pass`, and finally the database name with `--influx-db`.

```
./livedetect \
    --host 167.389.279.99 --port 4242 \
    --width 300 --height 300 \
    --detection \
    --create --repository /home/user/model/ \
    --confidence 0.1 --nclasses 601 --device-id 1 --gpu \
    --service openimage \
    -v INFO \
    --influx --influx-host http://192.168.90.90:8080 \
    --influx-user admin --influx-pass super-strong-password \
    --influx-db livedetect
```

## FAQ and Issues

- **Bounding box in the live feed are misplaced**: This is mostly due to your camera not resizing its output to fit the neural network input. In that case, look at the output dimension from the video output tab in your browser, and modify the `--width` and `--height` parameters accordingly. In practice this means the processing will be slower due to the processing of larger input frames

- **There's a delay in the video output**: This is a known issue that when the processing is much slower than frames input, a delay builds between the live input and output. Our current in-the-work solution to this is a specialized video input library named [libvnn](https://github.com/jolibrai/libvnn/)

- **Running DeepDetect on one machine and LiveDetect on another**: This is a common setup, typically running the DeepDetect container on a powerful GPU machine, and the LiveDetect binary somewhere near the camera, on a small board. Simply use the `--host` and `--port` parameters to fit your setup.

- **Running DeepDetect behind a reverse proxy**: If you use [DeepDetect Platform](https://deepdetect.com/platform/), your DeepDetect server will be accessible behind an nginx reverse proxy. You can use `--path` parameter to access it:

```
./livedetect \
  --host localhost \
  --port 1912 \
  --path /api/deepdetect \
  --mllib caffe \
  --width 300 --height 300 \
  --detection
```

- **Predict request on multiple services at once**: a json config file can be use to run multiple predict request on different services for each captured image. Here is an example, named `serviceConfig.json`:

```
[
  {
    "service": "detection_600",
    "parameters": {
      "input": {},
      "output": {
        "confidence_threshold": 0.5,
        "bbox": true
      },
      "mllib": {
        "gpu": true
      }
    }
  },
  {
    "service": "detectron_3k",
    "parameters": {
      "input": {},
      "output": {
        "confidence_threshold": 0.5,
        "bbox": true
      },
      "mllib": {
        "gpu": true
      }
    }
  }
]
```

Then use `livedetect` command with `--service-config` argument:

```
./livedetect \
  --host localhost \
  --port 8080 \
  --detection \
  --service-config ./serviceConfig.json \
  -v INFO
```
