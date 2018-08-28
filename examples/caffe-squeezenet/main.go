package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/milosgajdos83/ncs"
	"gocv.io/x/gocv"
)

var meanBGR = []float64{0.40787054 * 255.0, 0.45752458 * 255.0, 0.48109378 * 255.0}
var imgSize = image.Point{227, 227}

// meanCenter pre-preprocesses image so each of its layer pixels have zero mean
func meanCenter(img gocv.Mat, meanBGR []float64) gocv.Mat {
	r, c := img.Rows(), img.Cols()
	meanB, meanG, meanR := meanBGR[0], meanBGR[1], meanBGR[2]

	// create mean centered image layer by layer
	meanBMat := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(meanB, meanB, meanB, 0.0), r, c, gocv.MatTypeCV64F)
	meanGMat := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(meanB, meanG, meanG, 0.0), r, c, gocv.MatTypeCV64F)
	meanRMat := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(meanB, meanR, meanR, 0.0), r, c, gocv.MatTypeCV64F)
	meanMatImg := gocv.NewMat()
	gocv.Merge([]gocv.Mat{meanBMat, meanGMat, meanRMat}, &meanMatImg)
	defer meanMatImg.Close()

	floatImg := gocv.NewMat()
	defer floatImg.Close()
	img.ConvertTo(&floatImg, gocv.MatTypeCV64F)

	zeroMean := gocv.NewMat()
	gocv.Subtract(floatImg, meanMatImg, &zeroMean)

	return zeroMean
}

// readLabels reads labels file stored in labelsPath and returns it as a slice of strings
func readLabels(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// prepareImg preprocesses image for NCS
func prepareImg(img gocv.Mat) gocv.Mat {
	// resize the image
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Pt(227, 227), 0, 0, gocv.InterpolationDefault)
	// zero-mean centering
	zeroMeanImg := meanCenter(resized, meanBGR)
	// convert to FP32 for NCS
	fp32Image := gocv.NewMat()
	zeroMeanImg.ConvertTo(&fp32Image, gocv.MatTypeCV32F)

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
	dev, err := ncs.NewDevice(0)
	if err != nil {
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
	graph, err := ncs.NewGraph("SqueezenetGraph")
	if err != nil {
		return
	}
	defer graph.Destroy()
	log.Printf("NCS graph handle successfully created")

	graphFileName := "squeezenet_graph"
	graphData, err := ioutil.ReadFile(graphFileName)
	if err != nil {
		return
	}

	log.Printf("Attempting to allocate NCS graph")
	queue, err := graph.AllocateWithFifosOpts(dev, graphData, &ncs.FifoOpts{ncs.FifoHostWO, ncs.FifoFP32, 2}, &ncs.FifoOpts{ncs.FifoHostRO, ncs.FifoFP32, 2})
	if err != nil {
		return
	}
	defer queue.In.Destroy()
	defer queue.Out.Destroy()
	log.Printf("NCS Graph successfully allocated")

	log.Printf("Attempting to query INPUT FIFO options: %s", ncs.ROFifoElemDataSize)
	opts, err := queue.In.GetOption(ncs.ROFifoElemDataSize)
	if err != nil {
		return
	}
	data, err := ncs.ROFifoElemDataSize.Decode(opts)
	if err != nil {
		return
	}
	log.Printf("INPUT FIFO %s: %d", ncs.ROFifoElemDataSize, data.(uint))

	// digital image gymnastics
	imgPath := filepath.Join("nps_acoustic_guitar.png")
	log.Printf("Attempting to read image %s", imgPath)
	img := gocv.IMRead(imgPath, gocv.IMReadColor)

	ncsImg := prepareImg(img)
	log.Printf("Attempting to queue %s for inference", imgPath)
	err = graph.QueueInferenceWithFifoElem(queue, ncsImg.ToBytes(), nil)
	if err != nil {
		return
	}
	log.Printf("%s successfully queued for inference", imgPath)

	log.Printf("Attempting to read data from the OUTPUT FIFO queue")
	tensor, err := queue.Out.ReadElem()
	if err != nil {
		return
	}
	log.Printf("Read suceeded. Read %d bytes", len(tensor.Data))

	labelsPath := filepath.Join("squeeze_synset_words.txt")
	log.Printf("Reading labels file: %s", labelsPath)
	labels, err := readLabels(labelsPath)
	if err != nil {
		return
	}
	log.Printf("Read %d labels from %s", len(labels), labelsPath)

	// Decode the result
	var result [1000]float32
	buf := bytes.NewReader(tensor.Data)
	if err := binary.Read(buf, binary.LittleEndian, &result); err != nil {
		return
	}
	// find max and value
	max := result[0]
	idx := 0
	for i, val := range result {
		if max < val {
			max = val
			idx = i
		}
	}
	log.Printf("Prediction: %v, Probability: %v", labels[idx], max)
}
