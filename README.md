# LiveDetect

LiveDetect is a tool designed to process local video streams captured via camera, and execute machine learning models on each frame using [DeepDetect](https://github.com/jolibrain/deepdetect/). LiveDetect can be run on your desktop, on a Raspberry Pi 3 or on a Nvidia's GPU. 

Real-world use cases from DeepDetect customers with LiveDetect:

- Construction site safety and work monitoring

- Cars license plate OCR in parking lots

- Defect detection in manufactured precision pieces

![Example Traffic](example-traffic.gif)

## Set Up for Raspberry Pi 3

If you are using a **Raspberry Pi**, we made a step-by-step specific tutorial for the set up part [here](https://github.com/jolibrain/livedetect/wiki/Step-by-step-for-Raspberry-Pi-3). This tutorial is going to help you get your Raspberry Pi 3 ready. Once it's ready you can refer to the [examples section](https://github.com/jolibrain/livedetect#examples).

## Set Up for Deskstops and Nvidia's GPU

In order to use **Live Detect**, you need to have a DeepDetect instance running and then you need to build LiveDetect. 

### DeepDetect instance 

You need a **DeepDetect** instance running and accessible from the machine where you want to use LiveDetect.

Install the dependencies for LiveDetect:

- `sudo apt-get install libjpeg-dev`

First, create a `models` directory in your $HOME for the examples below to run without change. Then move to this new directory.

- **For Jetson ARM boards**

For Nvidia's Jetsons, you cannot use Docker and you should **build from source** as explained in the [quick start page](https://www.deepdetect.com/quickstart-server/?opts={%22os%22:%22ubuntu%22,%22source%22:%22build_source%22,%22compute%22:%22gpu%22,%22gpu%22:%22tx1%22,%22backend%22:[%22caffe%22],%22deepdetect%22:%22server%22}), using NCNN as a backend.

Once you are done with the quick start, you can directly [build LiveDetect from source](https://github.com/jolibrain/livedetect#buildfromsource).

- **Install Docker** (not for Jetson ARM boards)

We are going to need Docker, as it works very well with very little overhead. If [you already have Docker installed](https://www.digitalocean.com/community/questions/how-to-check-for-docker-installation), you can directly **start a DeepDetect container**.

To install Docker:

```
curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
sudo groupadd docker
sudo usermod -aG docker $USER
```

- **Start a DeepDetect container**

**If you are on your desktop**, you should make a CPU-only build, using Caffe:

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_cpu`

**If you are using a Nvidia's GPU**(not for Jetson ARM boards), you need to install [nvidia-docker](https://github.com/NVIDIA/nvidia-docker). Once it's installed, you can make a GPU build:

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_gpu`

### Build from source

You must go to the [releases page](https://github.com/jolibrain/livedetect/releases):

- If the binary for your system is available:

You should download the appropriate version of LiveDetect and then make it executable using:

`chmod +x livedetect-$VERSION` replacing `$VERSION` by the version of LiveDetect you downloaded.

It's ready! You can directly jump to the [examples section](https://github.com/jolibrain/livedetect#examples).

- If you couldn't find the LiveDetect version fit for your device, you should build from source by following the next steps.

#### Requirements

In order to build from source and use LiveDetect, you need Go.

To install Go, download the approriate version for your OS [here](https://golang.org/dl/).

When the download is complete, go to the directory the file has been downloaded to and extract it to install the go toolchain:

`tar -C $HOME -xzf go$VERSION.$OS-$ARCH.tar.gz`

A `go` directory will automatically be added to your `$HOME`. 

Finally, you should set the GoPATH. Add `$HOME/go/bin` to the PATH environment variable. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

`export PATH=$PATH:$HOME/go/bin`

**Note**: changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just refresh the file with `source $HOME/.profile`.

For more, refer to [this page](https://github.com/golang/go/wiki/SettingGOPATH).

#### LiveDetect build

Now that you have the required packages, go to your `$HOME` and fetch LiveDetect:

`go get -u github.com/jolibrain/livedetect`

This command should create a `gopath` directory 

Move in this directory, into to the LiveDetect file:

`cd ~/gopath/src/github.com/jolibrain/livedetect/`

From this directory, fetch all dependencies remaining:

`go get -u ./...`

Finally, build it!

`go build .`

You're ready for the following examples!

## Use-Cases

### Face detection + bounding boxes on web preview

This example starts LiveDetect and tells the DeepDetect instance located at 127.0.0.1:8080 to create a service named `face`.

This service takes 300x300 pictures for a detection process with the model located at the adress specified by `--init`, that have 2 classes.

The preview is displayed on port 8888, specified with `-P`.

V4L uses `--device-id` as the capture device (here device 0) and verbosity at the INFO level is triggered by `-v INFO`.

#### For Raspberry Pi

```
./livedetect-rpi3 \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/face/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_faces_ncnn_300x300.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```

#### For Desktop

marche pas
```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib caffe \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/face/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_faces_300x300.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```


```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib caffe \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/face/ \
    --init "https://www.deepdetect.com/models/init/desktop/images/detection/faces_512.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```


#### For Nvidia's GPU

meme modele que rpi3
```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/face/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_faces_ncnn_300x300.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    --gpu \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```

modele specifique
```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /opt/models/face/ \
    --init "https://www.deepdetect.com/models/init/embedded/images/detection/squeezenet_ssd_faces_ncnn.tar.gz " \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    --gpu \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```


### Detection + mask

This command starts LiveDetect and tells the DeepDetect instance located at localhost:4242 to create a service named `mask` for detection and mask preview.

It uses files under the `--repository` path with extensions files under the path specified by `-e`.

Mean values are specified with `-m`, 81 classes, using the GPU, service should be configured for images of size 416x608 and predictions should be considered only if their confidence is at least 0.1.

Finally, using `--select-classes`, we specify classes we want to be previewed and processed, using classes names with the `-c` flags.

```
./livedetect \
    --port 4242 --host localhost \
    --width 416 --height 608 \
    --detection \
    --create --repository /home/user/test_mask \
    --mask -e /home/user/test_mask/mask \
    -m "102.9801" -m "115.9465" -m "122.7717" \
    --confidence 0.1 --nclasses 81 \
    --service mask --device-id 1 \
    -P "0.0.0.0:8888" \
    --select-classes -c car -c person -c truck -c bike -c van
```
// caffe 2 dectron 

## Detection + InfluxDB

LiveDetect support InfluxDB, the categories detected and the probability for each category can be pushed to an InfluxDB database, here is the previous example, with the InfluxDB parameters.

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



