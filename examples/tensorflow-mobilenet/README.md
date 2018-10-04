# Mobilenet graph example

Example in this directory uses [MobileNetV2: Inverted Residuals and Linear Bottlenecks](https://arxiv.org/abs/1801.04381) [tensorflow](https://www.tensorflow.org/) model to classify an image of a giant panda:

<img src="./panda.jpg" alt="gian panda" width="227">

[image source](https://upload.wikimedia.org/wikipedia/commons/f/fe/Giant_Panda_in_Beijing_Zoo_1.JPG)

This example uses [mobilenet_v2_1.0_224](https://storage.googleapis.com/mobilenet_v2/checkpoints/mobilenet_v2_1.0_224.tgz). You can read more about other available MobilenetV2 models [here](https://github.com/tensorflow/models/tree/master/research/slim/nets/mobilenet).

If you would like to test other Mobilenet variants, you can use the `Makefile` available in this repo which allows you to download the model and compile it into Movidius graph file with a command like this:

```
make compile VERSION=v2 DEPTH=1.0 IMGSIZE=224
```

## Prerequisites

Install [GoCV](https://github.com/hybridgroup/gocv/#how-to-install.)

This example uses C/C++ NCSDK 2.0, so make sure you have it installed by following the instructions [here](https://movidius.github.io/ncsdk/install.html)

Note, since the 2.0 API does not seem to work properly on macOS, you won't be able to run this example on macOS. Everything works fine on Linux, in particular this example was tested on `Ubuntu 16.04`

## Running the example

Note, the example program contains hardcoded paths to the compiled Movidius graph file of the MobilenetV2, the image of giant pand and the MobilenetV2 labels.

You can run this example as follows:

```console
go run main.go
```

Result:

```console
2018/10/04 19:13:59 Attempting to create NCS device handle
2018/10/04 19:13:59 NCS device handle successfully created
2018/10/04 19:13:59 Attempting to open NCS device
2018/10/04 19:14:02 NCS device successfully opened
2018/10/04 19:14:02 Attempting to create NCS graph handle
2018/10/04 19:14:02 NCS graph handle successfully created
2018/10/04 19:14:02 Attempting to allocate NCS graph
2018/10/04 19:14:02 NCS Graph successfully allocated
2018/10/04 19:14:02 Attempting to read image panda.jpg
2018/10/04 19:14:02 Attempting to queue panda.jpg for inference
2018/10/04 19:14:02 panda.jpg successfully queued for inference
2018/10/04 19:14:02 Attempting to read data from the OUTPUT FIFO queue
2018/10/04 19:14:02 Read suceeded. Read 4004 bytes
2018/10/04 19:14:02 Reading labels file: imagenet_class_index.json
2018/10/04 19:14:02 Read 1000 labels from imagenet_class_index.json
2018/10/04 19:14:02 Prediction: [n02510455 giant_panda], Probability: 0.8408203
```
