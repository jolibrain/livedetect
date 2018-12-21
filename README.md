# LiveDetect

LiveDetect is a tool designed to process local video stream captured via camera, and execute machine learning models on every frame using [DeepDetect](https://github.com/jolibrain/deepdetect/).

![Example Traffic](example-traffic.gif)

## Quickstart

**NOTE:** this quickstart example can be followed perfectly on a Raspberry Pi 3B+, running [Raspbian](https://www.raspberrypi.org/downloads/raspbian/).

First, prepare your Pi, if you wish to use the *picamera* instead of an USB webcam activate it using:

- `sudo raspi-config`

Use the cursor keys to move to the camera option, and select 'enable'. On exiting  raspi-config, it will ask to reboot. The enable option will ensure that on reboot the correct GPU firmware will be running with the camera driver and tuning, and the GPU memory split is sufficient to allow the camera to acquire enough memory to run correctly.
Using the *raspi-config* tool, you can also activate SSH and expand the filesystem to the whole SD card, if you feel the need.

You will also need a decent amount of SWAP, check how much your system have:

- `free -m`

You should get this kind of output:

```
              total        used        free      shared  buff/cache   available
Mem:            927          31         828           0          67         845
Swap:          1023          30         993
```

If you notice that the total SWAP available (here 1023) is below 1000, follow the next steps.
To resize the SWAP to a decent amount, execute the following command:

- `sudo nano /etc/dphys-swapfile`

And change the **CONF_SWAPSIZE** value for 1024.

Finally, reload your SWAP setting:

- `sudo /sbin/dphys-swapfile setup`
- `sudo /sbin/dphys-swapfile swapoff`
- `sudo /sbin/dphys-swapfile swapon`

You're ready for the next step!
Install the dependencies for LiveDetect:

- `sudo apt install libjpeg-dev`

Download the binary for Raspberry Pi 3 system and make it executable:

- `wget https://github.com/jolibrain/livedetect/releases/download/v1.0.0/livedetect-rpi3`

- `chmod +x livedetect-rpi3`


**NOTE:** If you want to build LiveDetect by yourself, please refeer to the **Build** section of this README.

You need a **DeepDetect** instance running and accessible from the machine where you want to use LiveDetect.

If you want to run DeepDetect directly on a Raspberry Pi 3, here is a sample command to start a DeepDetect container with only NCNN as back-end (well suited for running directly on a Raspberry Pi).

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_ncnn_pi3`

For an AMD64 container, CPU-only, using Caffe:

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models beniz/deepdetect_cpu`

Or if you want a GPU build:

- `docker run -d -p 8080:8080 -v $HOME/models:/opt/models beniz/deepdetect_gpu`

**NOTE:** for the GPU container, you need to install [nvidia-docker](https://github.com/NVIDIA/nvidia-docker).

You're ready for the following examples!

## Examples

### Face detection + bounding boxes on web preview

```
./livedetect \
    --port 8080 \
    --host 127.0.0.1 \
    --mllib ncnn \
    --width 300 --height 300 \
    --detection \
    --create --repository /home/pi/face/ \
    --init "https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_faces_ncnn_300x300.tar.gz" \
    --confidence 0.3 \
    --device-id 0 \
    -v INFO \
    -P "0.0.0.0:8888" \
    --service face \
    --nclasses 2
```

This command start LiveDetect and tell the DeepDetect instance located at 127.0.0.1:8080 to create a service named `face-detect`.

This service takes 300x300 pictures for a detection process with the model located at the adress specified by `--init`, that have 2 classes.

The preview is displayed on port 8888, specified with `-P`.

V4L uses `--device-id` as the capture device (here device 0) and verbosity at the INFO level is triggered by `-v INFO`.

### Detection + mask

```
./LiveDetect \
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

This command start LiveDetect and tell the DeepDetect instance located at localhost:4242 to create a service named `mask` for detection and mask preview.

It uses files under the `--repository` path with extensions files under the path specified by `-e`.

Mean values are specified with `-m`, 81 classes, using the GPU, service should be configured for images of size 416x608 and predictions should be considered only if their confidence is at least 0.1.

Finally, using `--select-classes`, we specify classes we want to be previewed and processed, using classes names with the `-c` flags.

## Detection + InfluxDB

LiveDetect support InfluxDB, the categories detected and the probability for each category can be pushed to an InfluxDB database, here is the previous example, with the InfluxDB parameters.

```
./LiveDetect \
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

To use InfluxDB, you first need to pass the `--influx` argument, then you can specify the host with `--influx-host`, the credentials that should be used with `--influx-user` and `--influx-pass`, and finally the database name with `--influx-db`.

## Build from source

### Requirements

In order to build and use LiveDetect, you need Go.

Create your Go workspace, for example in your home directory:

`mkdir $HOME/go`

To install Go, download the approriate version for your OS [here](https://golang.org/dl/).

When the download is complete, extract it to install the go toolchain:

`tar -C $HOME -xzf go$VERSION.$OS-$ARCH.tar.gz`

Add `$HOME/go/bin` to the PATH environment variable. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

`export PATH=$PATH:$HOME/go/bin`

**Note**: changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just refresh the file with `source $HOME/.profile`.

Finally, set the GoPATH, for this, refeer to [this page](https://github.com/golang/go/wiki/SettingGOPATH).

### LiveDetect build

Now that you have the required packages fetch LiveDetect:

`go get -u github.com/jolibrain/livedetect`

Move in the directory containing the sources:

`cd $GOPATH/src/github.com/jolibrain/livedetect`

Fetch all dependencies remaining:

`go get -u ./...`

Finally, build it!

`go build .`
