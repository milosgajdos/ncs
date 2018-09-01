package ncs

// #cgo LDFLAGS: -lmvnc
/*
#include <ncs.h>
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

// GraphState defines states of a network graph
type GraphState int

const (
	// GraphCreated means the graph has been created, but it may not be initialized
	GraphCreated GraphState = iota
	// GraphAllocated means the graph has been initialized, and the graph has been allocated for a device
	GraphAllocated
	// GraphWaitingForInput means the graph is waiting for input.
	GraphWaitingForInput
	// GraphRunning means the graph is currently running an inference
	GraphRunning
)

// String implements fmt.Stringer interface for GraphState
func (gs GraphState) String() string {
	switch gs {
	case GraphCreated:
		return "GRAPH_CREATED"
	case GraphAllocated:
		return "GRAPH_ALLOCATED"
	case GraphWaitingForInput:
		return "GRAPH_WAITING_FOR_INPUT"
	case GraphRunning:
		return "GRAPH_RUNNING"
	default:
		return "GRAPH_UNKNOWN_STATE"
	}
}

// GraphOption defines graph options
// The options starting with RW are both gettable and settable
// The options starting with RO are only gettable
type GraphOption int

const (
	// ROGraphState is current state of the graph
	ROGraphState GraphOption = (1000 + iota)
	// ROGraphInferenceTime times taken per graph layer for the last inference in milliseconds
	ROGraphInferenceTime
	// ROGraphInputCount is number of inputs expected by the graph
	ROGraphInputCount
	// ROGraphOutputCount is he number of outputs expected by the graph.
	ROGraphOutputCount
	// ROGraphInputTensorDesc is an array of TensorDesc's, which describe the graph inputs in order
	ROGraphInputTensorDesc
	// ROGraphOutputTensorDesc is array of TensorDesc's, which describe the graph outputs in order
	ROGraphOutputTensorDesc
	// ROGraphDebugInfo provides more details when the result of a function call was StatusMyriadError
	ROGraphDebugInfo
	// ROGraphName is the name of the graph
	ROGraphName
	// ROGraphOptionClassLimit returns the highest option class supported
	ROGraphOptionClassLimit
	// ROGraphVersion is graph version
	ROGraphVersion
	// RWGraphExecutorsCount is not implemented yet
	RWGraphExecutorsCount
	// ROGraphInferenceTimeSize size of array for ROGraphInferenceTime option
	ROGraphInferenceTimeSize
)

// graphOptSize is a map which maps graph options to its native sizes
var graphOptSize = map[Option]uint{
	ROGraphState:             C.sizeof_int,
	ROGraphInferenceTime:     C.sizeof_float,
	ROGraphInputCount:        C.sizeof_int,
	ROGraphOutputCount:       C.sizeof_int,
	ROGraphInputTensorDesc:   C.sizeof_struct_ncTensorDescriptor_t,
	ROGraphOutputTensorDesc:  C.sizeof_struct_ncTensorDescriptor_t,
	ROGraphDebugInfo:         C.sizeof_char,
	ROGraphName:              C.sizeof_char,
	ROGraphOptionClassLimit:  C.sizeof_int,
	ROGraphVersion:           C.sizeof_char,
	RWGraphExecutorsCount:    C.sizeof_int,
	ROGraphInferenceTimeSize: C.sizeof_int,
}

// String implements fmt.Stringer interface for GraphOption
func (g GraphOption) String() string {
	switch g {
	case ROGraphState:
		return "GRAPH_STATE"
	case ROGraphInferenceTime:
		return "GRAPH_INFERENCE_TIME"
	case ROGraphInputCount:
		return "GRAPH_INPUT_COUNT"
	case ROGraphOutputCount:
		return "GRAPH_OUTPUT_COUNT"
	case ROGraphInputTensorDesc:
		return "GRAPH_INPUT_TENSOR_DESCRIPTORS"
	case ROGraphOutputTensorDesc:
		return "GRAPH_OUTPUT_TENSOR_DESCRIPTORS"
	case ROGraphDebugInfo:
		return "GRAPH_DEBUG_INFO"
	case ROGraphName:
		return "GRAPH_NAME"
	case ROGraphOptionClassLimit:
		return "GRAPH_OPTION_CLASS_LIMIT"
	case ROGraphVersion:
		return "GRAPH_VERSION"
	case RWGraphExecutorsCount:
		return "RW_GRAPH_EXECUTORS_LIMIT"
	case ROGraphInferenceTimeSize:
		return "GRAPH_INFERENCE_TIME_SIZE"
	default:
		return "GRAPH_UNKNOWN_OPTION"
	}
}

// Value returns option value as integer
func (g GraphOption) Value() int {
	return int(g)
}

// Decode decodes options data encoded in raw bytes and returns it in its native type.
// The returned data then can be asserted into its native type.
// If the data contains more than one element you need to specify the number of expected elements via count.
// It returns error if the data fails to be decoded into the option native type.
func (g GraphOption) Decode(data []byte, count int) (interface{}, error) {
	buf := bytes.NewReader(data)

	switch g {
	case ROGraphState,
		ROGraphInputCount,
		ROGraphOutputCount,
		ROGraphOptionClassLimit,
		RWGraphExecutorsCount,
		ROGraphInferenceTimeSize:

		var val uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		// this is safe type cast as we know val is a positive integer
		return uint(val), nil

	case ROGraphInferenceTime:
		val := make([]float32, count)
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return val[:], nil

	case ROGraphVersion:

		var val [2]uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return val[:], nil

	case ROGraphDebugInfo,
		ROGraphName:

		return string(data), nil

	case ROGraphInputTensorDesc,
		ROGraphOutputTensorDesc:
		vals := make([]struct {
			BatchSize uint32
			Channels  uint32
			Width     uint32
			Height    uint32
			Size      uint32
			CStride   uint32
			WStride   uint32
			HStride   uint32
			DataType  int32
		}, count)

		if err := binary.Read(buf, binary.LittleEndian, &vals); err != nil {
			return nil, err
		}

		tensorDescs := make([]TensorDesc, count)
		for i, val := range vals {
			td := TensorDesc{
				BatchSize: uint(val.BatchSize),
				Channels:  uint(val.Channels),
				Width:     uint(val.Width),
				Height:    uint(val.Height),
				Size:      uint(val.Size),
				CStride:   uint(val.CStride),
				WStride:   uint(val.WStride),
				HStride:   uint(val.HStride),
				DataType:  FifoDataType(val.DataType),
			}
			tensorDescs[i] = td
		}

		return tensorDescs[:], nil

	default:
		return nil, fmt.Errorf("Unable to decode graph option data: %s", g)
	}

	return nil, nil
}

// Graph is NCSDK neural network graph
type Graph struct {
	name   string
	handle unsafe.Pointer
	device *Device
}

// NewGraph creates new Graph with given name and returns it
// It returns error if it fails to create new graph
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphCreate.html
func NewGraph(name string) (*Graph, error) {
	var handle unsafe.Pointer

	_name := C.CString(name)
	defer C.free(unsafe.Pointer(_name))

	s := C.ncs_GraphCreate(_name, &handle)

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to create new graph: %s", Status(s))
	}

	return &Graph{name: name, handle: handle}, nil
}

// Allocate allocates a graph on NCS device. This function sends graphData to NCS device. It does not allocate input or output FIFO queues. You have to either allocate them separately or use either AllocateWithFifosDefault() or AllocateWithFifosOpts() functions whcih conveniently create and allocate the FIFO queues.
// It returns error if it fails to allocate the graph on the device
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphAllocate.html
func (g *Graph) Allocate(d *Device, graphData []byte) error {
	s := C.ncs_GraphAllocate(d.handle, g.handle, unsafe.Pointer(&graphData[0]), C.uint(len(graphData)))

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to allocate new graph: %s", Status(s))
	}

	g.device = d

	return nil
}

// AllocateWithFifosDefault allocates a graph and creates and allocates FIFO queues with default parameters for inference. Both FIFOs have FifoDataType set to FifoFP32. Inbound FIFO queue is initialized with FifoHostWO type and outbound FIFO queue with FifoHostRO type. It returns FifoQueue or error if it fails to allocate the graph.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphAllocateWithFifos.html
func (g *Graph) AllocateWithFifosDefault(d *Device, graphData []byte) (*FifoQueue, error) {
	return g.AllocateWithFifosOpts(d, graphData, &FifoOpts{FifoHostWO, FifoFP32, 2}, &FifoOpts{FifoHostRO, FifoFP32, 2})
}

// AllocateWithFifosOpts allocates a graph and creates and allocates FIFO queues for inference. This function is similar to AllocateWithFifosDefault, but rather than initializing FIFOs with default values it accepts parameters that allow to specify FIFO queue parameters
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphAllocateWithFifosEx.html
func (g *Graph) AllocateWithFifosOpts(d *Device, graphData []byte, inOpts *FifoOpts, outOpts *FifoOpts) (*FifoQueue, error) {
	var inHandle, outHandle unsafe.Pointer

	s := C.ncs_GraphAllocateWithFifosEx(d.handle,
		g.handle, unsafe.Pointer(&graphData[0]), C.uint(len(graphData)),
		&inHandle, C.ncFifoType(inOpts.Type), C.int(inOpts.NumElem), C.ncFifoDataType(inOpts.DataType),
		&outHandle, C.ncFifoType(outOpts.Type), C.int(outOpts.NumElem), C.ncFifoDataType(outOpts.DataType))

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to allocate graph with FIFOs: %s", Status(s))
	}

	g.device = d

	return &FifoQueue{
		In:  &Fifo{handle: inHandle, device: d},
		Out: &Fifo{handle: outHandle, device: d},
	}, nil
}

// QueueInference queues data for inference to be processed by a graph with specified input and output FIFOs
// If it fails to queue the data tensor it returns error
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphQueueInference.html
func (g *Graph) QueueInference(f *FifoQueue) error {
	s := C.ncs_GraphQueueInference(g.handle, &f.In.handle, C.uint(1), &f.Out.handle, C.uint(1))

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to queue inference: %s", Status(s))
	}

	return nil
}

// QueueInferenceWithFifoElem writes an element to a FIFO, usually an input tensor for inference, and queues an inference to be processed by a graph. This is a convenient way to write an input tensor and queue an inference in one call
// If it fails to queue the data tensor it returns error
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphQueueInferenceWithFifoElem.html
func (g *Graph) QueueInferenceWithFifoElem(f *FifoQueue, data []byte, metaData interface{}) error {
	dataLen := C.uint(len(data))

	s := C.ncs_GraphQueueInferenceWithFifoElem(g.handle, f.In.handle, f.Out.handle, unsafe.Pointer(&data[0]), &dataLen, unsafe.Pointer(&metaData))

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to queue inference: %s", Status(s))
	}

	return nil
}

// GetOption queries the value of an option for a graph and returns it encoded in a byte slice
// It returns error if it failed to retrieve the option value
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphGetOption.html
func (g *Graph) GetOption(opt GraphOption) ([]byte, error) {
	if opt == RWGraphExecutorsCount {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	var data unsafe.Pointer
	var dataLen C.uint

	s := C.ncs_GraphGetOption(g.handle, C.int(opt), data, &dataLen)

	if Status(s) == StatusInvalidDataLength {
		return g.GetOptionWithByteSize(opt, graphOptSize[opt]*uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option: %s", opt, Status(s))
}

// GetOptionsWithSize queries NCS grapg options and returns it encoded in a byte slice of size elements.
// This function is similar to GetOption(), however as opposed to GetOption() which first queries the NCS device for the size of the requested options, it attempts to request the options data by specifying its size in raw bytes explicitly, hence it returns the queried options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphGetOption.html
func (g *Graph) GetOptionWithByteSize(opt GraphOption, size uint) ([]byte, error) {
	if opt == RWGraphExecutorsCount {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("graph", g.handle, opt, size)
}

// Destroy destroys NCS graph handle and frees associated resources.
// This function must be called for every graph that was initialized with NewGraph().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphDestroy.html
func (g *Graph) Destroy() error {
	s := C.ncs_GraphDestroy(&g.handle)

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to destroy graph: %s", Status(s))
	}

	return nil
}
