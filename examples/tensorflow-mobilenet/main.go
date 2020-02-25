package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/milosgajdos/ncs"
	"gocv.io/x/gocv"
)

// readLabels reads labels file stored in labelsPath and returns it as a map
// the maps keys are the label numbers; the maps value is a slice that contains label metadata
func readLabels(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var labels map[string][]string
	dec := json.NewDecoder(file)
	if err := dec.Decode(&labels); err != nil {
		return nil, err
	}

	return labels, nil
}

// prepareImg preprocesses image for NCS
func prepareImg(img gocv.Mat) gocv.Mat {
	// resize the image
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Pt(224, 224), 0, 0, gocv.InterpolationDefault)

	fp32Image := gocv.NewMat()
	resized.ConvertTo(&fp32Image, gocv.MatTypeCV32F)

	fp32Image.DivideFloat(128.0)
	fp32Image.SubtractFloat(1.0)

	return fp32Image
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}()
	log.Printf("Attempting to create NCS device handle")
	dev, e := ncs.NewDevice(0)
	if e != nil {
		err = e
		return
	}
	defer dev.Destroy()
	log.Printf("NCS device handle successfully created")

	log.Printf("Attempting to open NCS device")
	err = dev.Open()
	if err != nil {
		return
	}
	defer dev.Close()
	log.Printf("NCS device successfully opened")

	log.Printf("Attempting to create NCS graph handle")
	graph, e := ncs.NewGraph("MobilenetGraph")
	if e != nil {
		err = e
		return
	}
	defer graph.Destroy()
	log.Printf("NCS graph handle successfully created")

	graphFileName := "mobilenet_graph"
	graphData, e := ioutil.ReadFile(graphFileName)
	if e != nil {
		err = e
		return
	}

	log.Printf("Attempting to allocate NCS graph")
	queue, e := graph.AllocateWithFifosOpts(dev, graphData,
		&ncs.FifoOpts{ncs.FifoHostWO, ncs.FifoFP32, 2},
		&ncs.FifoOpts{ncs.FifoHostRO, ncs.FifoFP32, 2})
	if e != nil {
		err = e
		return
	}
	defer queue.In.Destroy()
	defer queue.Out.Destroy()
	log.Printf("NCS Graph successfully allocated")

	// digital image gymnastics
	imgPath := filepath.Join("panda.jpg")
	log.Printf("Attempting to read image %s", imgPath)
	img := gocv.IMRead(imgPath, gocv.IMReadColor)

	ncsImg := prepareImg(img)
	log.Printf("Attempting to queue %s for inference", imgPath)
	err = graph.QueueInferenceWithFifoElem(queue, ncsImg.ToBytes(), nil)
	if err != nil {
		return
	}
	log.Printf("%s successfully queued for inference", imgPath)

	log.Printf("Attempting to read data from NCS")
	tensor, e := queue.Out.ReadElem()
	if e != nil {
		err = e
		return
	}
	log.Printf("Read suceeded. Read %d bytes", len(tensor.Data))

	labelsPath := filepath.Join("imagenet_class_index.json")
	log.Printf("Reading labels file: %s", labelsPath)
	labels, e := readLabels(labelsPath)
	if e != nil {
		err = e
		return
	}
	log.Printf("Read %d labels from %s", len(labels), labelsPath)

	// Decode the result
	var result [1000]float32
	buf := bytes.NewReader(tensor.Data)
	if e := binary.Read(buf, binary.LittleEndian, &result); e != nil {
		err = e
		return
	}

	// find max and label
	max := result[0]
	idx := 0
	for i, val := range result {
		if max < val {
			max = val
			idx = i
		}
	}
	// need to subtract 1 from idx as Mobilenet indexes from 1, our labels index from 0
	// https://github.com/mldbai/tensorflow-models/blob/master/inception/inception/data/imagenet_lsvrc_2015_synsets.txt
	log.Printf("Prediction: %v, Probability: %v", labels[strconv.Itoa(idx-1)], max)
}
