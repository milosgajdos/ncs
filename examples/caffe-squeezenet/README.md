# SqueezeNet graph example

Example in this directory uses [SqueezeNet](https://arxiv.org/abs/1602.07360) caffe model to classify the following image of acoustic guitar:

<img src="./nps_acoustic_guitar.png" alt="acoustic guitar" width="227">


## Prerequisites

This example uses C/C++ NCSDK 2.0, so make sure you have it installed by following the instructions [here](https://movidius.github.io/ncsdk/install.html)

## Running the example

You can run this example as follows

```console
go run main.go
```

Result:

```console
2018/08/28 13:48:24 Attempting to create NCS device handle
2018/08/28 13:48:24 NCS device handle successfully created
2018/08/28 13:48:24 Attempting to open NCS device
2018/08/28 13:48:28 NCS device successfully opened
2018/08/28 13:48:28 Attempting to create NCS graph handle
2018/08/28 13:48:28 NCS graph handle successfully created
2018/08/28 13:48:28 Attempting to allocate NCS graph
2018/08/28 13:48:28 NCS Graph successfully allocated
2018/08/28 13:48:28 Attempting to query INPUT FIFO options: RO_FIFO_ELEM_DATA_SIZE
2018/08/28 13:48:28 INPUT FIFO RO_FIFO_ELEM_DATA_SIZE: 618348
2018/08/28 13:48:28 Attempting to read image nps_acoustic_guitar.png
2018/08/28 13:48:28 Attempting to queue nps_acoustic_guitar.png for inference
2018/08/28 13:48:28 nps_acoustic_guitar.png successfully queued for inference
2018/08/28 13:48:28 Attempting to read data from the OUTPUT FIFO queue
2018/08/28 13:48:28 Read suceeded. Read 4000 bytes
2018/08/28 13:48:28 Reading labels file: squeeze_synset_words.txt
2018/08/28 13:48:28 Read 1000 labels from squeeze_synset_words.txt
2018/08/28 13:48:28 Prediction: n02676566 acoustic guitar, Probability: 0.99316406
```
