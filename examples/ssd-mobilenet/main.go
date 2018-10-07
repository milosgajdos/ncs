package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/milosgajdos83/ncs"
	"gocv.io/x/gocv"
)

// checks if the number is finite i.e. not NaN or in (-Inf, Inf) interval
func isFinite(f float64) bool {
	if math.IsNaN(f) {
		return false
	}

	if math.IsInf(f, 1) {
		return false
	}

	if math.IsInf(f, -1) {
		return false
	}

	return true
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
	gocv.Resize(img, &resized, image.Pt(300, 300), 0, 0, gocv.InterpolationDefault)

	floatImg := gocv.NewMat()
	defer floatImg.Close()
	resized.ConvertTo(&floatImg, gocv.MatTypeCV64F)

	fp32Image := gocv.NewMat()
	floatImg.ConvertTo(&fp32Image, gocv.MatTypeCV32F)

	fp32Image.SubtractFloat(127.5)
	fp32Image.MultiplyFloat(0.007843)

	return fp32Image
}

// drawBoxes draws the boxes and labels of all detected objects into the original image
func drawBoxesAndLabels(img, result gocv.Mat, labels []string) {
	boxes := int(result.GetFloatAt(0, 0))
	rows, cols := img.Rows(), img.Cols()
	for i := 0; i < boxes; i++ {
		idx := 7 + i*7
		if !(isFinite(float64(result.GetFloatAt(0, idx))) &&
			isFinite(float64(result.GetFloatAt(0, idx+1))) &&
			isFinite(float64(result.GetFloatAt(0, idx+2))) &&
			isFinite(float64(result.GetFloatAt(0, idx+3))) &&
			isFinite(float64(result.GetFloatAt(0, idx+4))) &&
			isFinite(float64(result.GetFloatAt(0, idx+5))) &&
			isFinite(float64(result.GetFloatAt(0, idx+6)))) {
			continue
		}

		x1 := math.Max(0, float64(result.GetFloatAt(0, idx+3))*float64(rows))
		y1 := math.Max(0, float64(result.GetFloatAt(0, idx+4))*float64(cols))
		x2 := math.Min(float64(rows), float64(result.GetFloatAt(0, idx+5))*float64(rows))
		y2 := math.Min(float64(cols), float64(result.GetFloatAt(0, idx+6))*float64(cols))

		classID := labels[int(result.GetFloatAt(0, idx+1))]
		confidence := result.GetFloatAt(0, idx+2) * 100.0

		log.Printf("Box at: %d: ClassID: %s Confidence: %.2f, Top Left: (%d, %d) Bottom Right: (%d, %d)",
			i, classID, confidence, int(x1), int(y1), int(x2), int(y2))

		gocv.Rectangle(&img, image.Rect(int(x1), int(y1), int(x2), int(y2)), color.RGBA{0, 0, 255, 0}, 2)

		label := fmt.Sprintf("%s: %.2f", classID, confidence)
		labelBgColor := color.RGBA{125, 175, 75, 0}
		labelTxtColor := color.RGBA{255, 255, 255, 0}
		labelSize := gocv.GetTextSize(label, gocv.FontHersheySimplex, 0.5, 1)

		lx1 := int(x1)
		ly1 := int(y1) - labelSize.Y
		if ly1 < 1 {
			ly1 = 1
		}
		lx2 := lx1 + labelSize.X
		ly2 := ly1 + labelSize.Y

		gocv.Rectangle(&img, image.Rect(int(lx1), int(ly1), int(lx2), int(ly2)), labelBgColor, -1)
		gocv.PutText(&img, label, image.Pt(lx1, ly2), gocv.FontHersheySimplex, 0.5, labelTxtColor, 1)
	}
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
	if err != nil {
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
	graph, e := ncs.NewGraph("SSDMobilenetGraph")
	if e != nil {
		err = e
		return
	}
	defer graph.Destroy()
	log.Printf("NCS graph handle successfully created")

	graphFileName := "ssd_mobilenet_graph"
	graphData, e := ioutil.ReadFile(graphFileName)
	if e != nil {
		err = e
		return
	}

	log.Printf("Attempting to allocate NCS graph")
	queue, err := graph.AllocateWithFifosOpts(dev, graphData,
		&ncs.FifoOpts{ncs.FifoHostWO, ncs.FifoFP16, 2},
		&ncs.FifoOpts{ncs.FifoHostRO, ncs.FifoFP16, 2})
	if e != nil {
		err = e
		return
	}
	defer queue.In.Destroy()
	defer queue.Out.Destroy()
	log.Printf("NCS Graph successfully allocated")

	// digital image gymnastics
	imgPath := filepath.Join("nps_chair.png")
	log.Printf("Attempting to read image %s", imgPath)
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	// need to submit FP16 image
	fp32Image := prepareImg(img)
	ncsImg := fp32Image.ConvertFp16()
	log.Printf("Attempting to queue %s for inference", imgPath)
	err = graph.QueueInferenceWithFifoElem(queue, ncsImg.ToBytes(), nil)
	if err != nil {
		return
	}
	log.Printf("%s successfully queued for inference", imgPath)

	log.Printf("Attempting to read data from NCS")
	tensor, err := queue.Out.ReadElem()
	if e != nil {
		err = e
		return
	}
	log.Printf("Read suceeded. Read %d bytes", len(tensor.Data))

	labelsPath := filepath.Join("labels.txt")
	log.Printf("Reading labels file: %s", labelsPath)
	labels, err := readLabels(labelsPath)
	if e != nil {
		err = e
		return
	}
	log.Printf("Read %d labels from %s", len(labels), labelsPath)

	// Result is returned as 32bit float, but we are data into 16bit floats hence we
	// need only half the size of the original 32bit float result data buffer
	fp16Result, err := gocv.NewMatFromBytes(1, len(tensor.Data)/2, gocv.MatTypeCV16S, tensor.Data)
	if e != nil {
		err = e
		return
	}

	result := fp16Result.ConvertFp16()
	log.Printf("Detected boxes: %d", int(result.GetFloatAt(0, 0)))
	drawBoxesAndLabels(img, result, labels)

	resultPath := filepath.Join("result.png")
	if !gocv.IMWrite(resultPath, img) {
		err = fmt.Errorf("Failed to save image %s", resultPath)
		return
	}
}
