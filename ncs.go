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

// getOption is a utility function which lets you query various NCS options
func getOption(resource string, handle unsafe.Pointer, option int, size uint) ([]byte, error) {
	// allocate buffer for options data
	data := C.malloc(C.sizeof_char * C.ulong(size))
	defer C.free(unsafe.Pointer(data))
	dataLen := C.uint(size)

	// NCCS API status code
	var s C.int

	switch resource {
	case "device":
		s = C.ncs_DeviceGetOption(handle, C.int(option), data, &dataLen)
	case "graph":
		s = C.ncs_GraphGetOption(handle, C.int(option), data, &dataLen)
	case "fifo":
		s = C.ncs_FifoGetOption(handle, C.int(option), data, &dataLen)
	default:
		return nil, fmt.Errorf("Unknown resource: %s", resource)
	}

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to get %s option: %s", resource, Status(s))
	}

	return C.GoBytes(unsafe.Pointer(data), C.int(size)), nil
}

const (
	// ThermalBufferSize the size of the temperature buffer as returned when querying device
	ThermalBufferSize = 100
	// MaxNameSize is the maximum lengrth of device or graph name size
	MaxNameSize = 28
)

// Status is the NCSDK API status code as returned by most API calls.
// It usually reports the status of the neural compute stick.
type Status int

const (
	// StatusOK means the API function call worked as expected
	StatusOK Status = -iota
	// StatusBusy means device is busy, retry later.
	StatusBusy
	// StatusError means an unexpected error was encountered during the API function call.
	StatusError
	// StatusOutOfMemory means the host is out of memory.
	StatusOutOfMemory
	// StatusDeviceNotFound means no device has been found at the given index or name.
	StatusDeviceNotFound
	// StatusInvalidParameters means at least one of the given parameters is wrong.
	StatusInvalidParameters
	// StatusTimeout means there was a timeout in the communication with the device.
	StatusTimeout
	// StatusCmdNotFound means the file to boot the device was not found.
	StatusCmdNotFound
	// StatusNotAllocated means the graph or fifo has not been allocated.
	StatusNotAllocated
	// StatusUnauthorized means an unauthorized operation has been attempted.
	StatusUnauthorized
	// StatusUnsupportedGraphFile means the graph file version is not supported.
	StatusUnsupportedGraphFile
	// StatusUnsupportedConfigFile is reserved for future use.
	StatusUnsupportedConfigFile
	// StatusUnsupportedFeature means the operation used a feature unsupported by this firmware version.
	StatusUnsupportedFeature
	// StatusMyriadError when an error has been reported by device, use MVNC_DEBUG_INFO.
	StatusMyriadError
	// StatusInvalidDataLength means an invalid data length has been passed when getting or setting an option
	StatusInvalidDataLength
	// StatusInvalidHandle means an invalid handle has been passed to a function
	StatusInvalidHandle
)

// String method to satisfy fmt.Stringer interface
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "STATUS_OK"
	case StatusBusy:
		return "DEVICE_BUSY"
	case StatusError:
		return "UNEXPECTED_ERROR"
	case StatusOutOfMemory:
		return "HOST_OUT_OF_MEMORY"
	case StatusDeviceNotFound:
		return "DEVICE_NOT_FOUND"
	case StatusInvalidParameters:
		return "INVALID_PARAMETERS"
	case StatusTimeout:
		return "TIMEOUT"
	case StatusCmdNotFound:
		return "BOOTLOADER_NOT_FOUND"
	case StatusNotAllocated:
		return "UNALLOCATED_RESOURCE"
	case StatusUnauthorized:
		return "UNAUTHORIZED_OPERATION"
	case StatusUnsupportedGraphFile:
		return "UNSUPPORTED_GRAPH_FILE"
	case StatusUnsupportedConfigFile:
		return "UNSUPPORTED_CONFIGURATION"
	case StatusUnsupportedFeature:
		return "UNSUPPORTED_FEATURE"
	case StatusMyriadError:
		return "MOVIDIUS_VPU_ERROR"
	case StatusInvalidDataLength:
		return "INVALID_OPTION_LENGTH"
	case StatusInvalidHandle:
		return "INVALID_HANDLE"
	default:
		return "UNKNOWN_STATUS"
	}
}

// DeviceOption defines NCS options
type DeviceOption int

const (
	// RODeviceThermalStats allows to query device temperatures in degrees Celsius
	// This option returns []float64 array of maxi temperatures for the last ThermalBufferSize seconds.
	RODeviceThermalStats DeviceOption = (2000 + iota)
	// RODeviceThermalThrottling allows to query temperature throttling level
	RODeviceThermalThrottle
	// RODeviceState qieries state of the device
	RODeviceState
	// RODeviceMemoryUsed allows to query current memory in use on the device in bytes.
	// Returned value must be divided by RODeviceMemorySize if the percentage of memory needs to be computed
	RODeviceMemoryUsed
	// RODeviceMemorySize queries total memory available on the device in bytes
	RODeviceMemorySize
	// RODeviceMaxFifoCount queries maximum number of FIFOs that can be allocated for the device
	RODeviceMaxFifoCount
	// RODeviceAllocatedFifoCount queries number of FIFOs currently allocated for the device
	RODeviceAllocatedFifoCount
	// RODeviceMaxMaxGraphCount queries the maximum number of graphs that can be allocated for the device
	RODeviceMaxGraphCount
	// RODeviceAllocatedGraphCount queries the number of graphs currently allocated for the device
	RODeviceAllocatedGraphCount
	// RODeviceClassLimit queries the highest option class supported
	RODeviceClassLimit
	// RODeviceFirmwareVersion queries queries the version of the firmware currently running on the device
	RODeviceFirmwareVersion
	// RODeviceDebugInfo queries more detailed info when the result of a function call was StatusMyriadError
	RODeviceDebugInfo
	// RODeviceMVTensorVersion queries the version of the mvtensor library that was linked with the API
	RODeviceMVTensorVersion
	// RODeviceName queries the internal name of the device
	RODeviceName
	// RODeviceMaxExecutors is reserved for future use
	RODeviceMaxExecutors
	// RODeviceHWVersion queries the hardware version of the device
	RODeviceHWVersion
)

// String implements fmt.Stringer interface for DeviceOption
func (do DeviceOption) String() string {
	switch do {
	case RODeviceThermalStats:
		return "RO_DEVICE_THERMAL_STATE"
	case RODeviceThermalThrottle:
		return "RO_DEVICE_THERMAL_THROTTLE"
	case RODeviceState:
		return "RO_DEVICE_STATE"
	case RODeviceMemoryUsed:
		return "RO_DEVICE_MEMORY_USED"
	case RODeviceMemorySize:
		return "RO_DEVICE_MEMORY_SIZE"
	case RODeviceMaxFifoCount:
		return "RO_DEVICE_MAX_FIFO_COUNT"
	case RODeviceAllocatedFifoCount:
		return "RO_DEVICE_ALLOCATED_FIFO_COUNT"
	case RODeviceMaxGraphCount:
		return "RO_DEVICE_MAX_GRAPH_COUNT"
	case RODeviceAllocatedGraphCount:
		return "RO_DEVICE_ALLOCATED_GRAPH_COUNT"
	case RODeviceClassLimit:
		return "RO_DEVICE_CLASS_LIMIT"
	case RODeviceFirmwareVersion:
		return "RO_DEVICE_FIRMWARE_VERSION"
	case RODeviceDebugInfo:
		return "RO_DEVICE_DEBUG_INFO"
	case RODeviceMVTensorVersion:
		return "RO_DEVICE_MVTENSOR_VERSION"
	case RODeviceName:
		return "RO_DEVICE_NAME"
	case RODeviceMaxExecutors:
		return "RO_DEVICE_MAX_EXECUTORS"
	case RODeviceHWVersion:
		return "RO_DEVICE_HW_VERSION"
	default:
		return "DEVICE_UNKNOWN_OPTION"
	}
}

// DeviceState defines NCS device status
type DeviceState int

const (
	// DeviceCreated means device has been created
	DeviceCreated DeviceState = iota
	// DeviceOpened means device has been opened
	DeviceOpened
	// DeviceClosed means device has been closed
	DeviceClosed
)

// String implements fmt.Stringer interface for DeviceState
func (ds DeviceState) String() string {
	switch ds {
	case DeviceCreated:
		return "DEVICE_CREATED"
	case DeviceOpened:
		return "DEVICE_OPENED"
	case DeviceClosed:
		return "DEVICE_CLOSED"
	default:
		return "DEVICE_UNKNOWN_STATUS"
	}
}

// Device is Neural Compute Stick (NCS) device
type Device struct {
	handle unsafe.Pointer
}

// NewDevice creates new NCS device handle and returns it
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceCreate.html
func NewDevice(index int) (*Device, error) {
	var handle unsafe.Pointer

	s := C.ncs_DeviceCreate(C.int(index), &handle)

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to create new device: %s", Status(s))
	}

	return &Device{handle: handle}, nil
}

// Open initializes NCS device and opens device communication channel
// It returns error if it fails to open the device.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceOpen.html
func (d *Device) Open() error {
	s := C.ncs_DeviceOpen(d.handle)

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to open device: %s", Status(s))
	}

	return nil
}

// GetOption queries the value of an option for the device and returns it in a byte slice
// It returns error if it failed to retrieve the option value
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceGetOption.html
func (d *Device) GetOption(opt DeviceOption) ([]byte, error) {
	if opt == RODeviceMaxExecutors || opt == RODeviceDebugInfo {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	var data unsafe.Pointer
	var dataLen C.uint

	s := C.ncs_DeviceGetOption(d.handle, C.int(opt), data, &dataLen)

	if Status(s) == StatusInvalidDataLength {
		return d.GetOptionWithSize(opt, uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option", opt)
}

// GetOptionsWithSize queries device options and returns it encoded in a byte slice of the same size as requested if possible. This function is similar to GetOption(), however as opposed to figuring out the byte size of the queried options it attempts to request the options data by specifying its size explicitly. Because we specify the options data size explicitly this function returns the options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceGetOption.html
func (d *Device) GetOptionWithSize(opt DeviceOption, size uint) ([]byte, error) {
	if opt == RODeviceMaxExecutors || opt == RODeviceDebugInfo {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("device", d.handle, int(opt), size)
}

// Close closes the communication channel with NCS device.
// It returns error if it fails to close the communication channel.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceClose.html
func (d *Device) Close() error {
	s := C.ncs_DeviceClose(d.handle)

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to close device: %s", Status(s))
	}

	return nil
}

// Destroy destroys NCS device handle and frees associated resources.
// This function must be called for every device that was initialized with NewDevice().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceDestroy.html
func (d *Device) Destroy() error {
	s := C.ncs_DeviceDestroy(&d.handle)

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to destroy device: %s", Status(s))
	}

	return nil
}

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
		return "GRAPH_INPUT_TENSOR_DESCRIPTION"
	case ROGraphOutputTensorDesc:
		return "GRAPH_OUTPUT_TENSOR_DESCRIPTION"
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
		return g.GetOptionWithSize(opt, uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option", opt)
}

// GetOptionsWithSize queries graph options and returns it encoded in a byte slice of the same size as requested if possible. This function is similar to GetOption(), however as opposed to figuring out the byte size of the queried options it attempts to request the options data by specifying its size explicitly. Because we specify the options data size explicitly this function returns the options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphGetOption.html
func (g *Graph) GetOptionWithSize(opt GraphOption, size uint) ([]byte, error) {
	if opt == RWGraphExecutorsCount {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("graph", g.handle, int(opt), size)
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

// TensorDesc describes graph inputs and outputs
type TensorDesc struct {
	// BatchSize contains number of elements
	BatchSize uint
	// Channels contains number of channels (when dealing with digital images)
	Channels uint
	// Width is data width (i.e. number of matrix columns)
	Width uint
	// Height is data height (i.e. number of matrix rows)
	Height uint
	// Size is total data size in tensor
	Size uint
	// CStride is channel stride (Stride in the channels' dimension)
	CStride uint
	// WStride is width stride (Stride in the horizontal dimension)
	WStride uint
	// HStride is height stride (Stride in the vertical dimension)
	HStride uint
	// DataType is data type of the tensor
	DataType FifoDataType
}

// Tensor is graph tensor as returned from NCS
type Tensor struct {
	// Data contains raw tensor data
	Data []byte
	// MetaData contains tensor metadata
	MetaData interface{}
}

// FifoQueue is a FIFO queue used for NCS inference
type FifoQueue struct {
	// In is an inbound queue
	In *Fifo
	// Out is an outbound queue
	Out *Fifo
}

// FifoType defines FIFO access types
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoType_t.html
type FifoType int

const (
	// FifoHostRO allows Read Only API access and Read-Write Graph access
	FifoHostRO FifoType = iota
	// FifoHostWO allows Write Only API acess and Read Only Graph access
	FifoHostWO
)

// FifoDataType defines possible data types for FIFOs
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoDataType_t.html
type FifoDataType int

const (
	// FifoFP16 data is in half precision (16 bit) floating point format (FP16).
	FifoFP16 FifoDataType = iota
	// FifoFP32 data is in full precision (32 bit) floating point format (FP32)
	FifoFP32
)

// String implements fmt.Stringer interface for FifoDataType
func (fd FifoDataType) String() string {
	switch fd {
	case FifoFP16:
		return "FIFO_FLOAT_16"
	case FifoFP32:
		return "FIFO_FLOAT_32"
	default:
		return "FFIO_UNKNOWN_DATA_TYPE"
	}
}

// FifoState is state of FIFO
type FifoState int

const (
	// FifoCreated means FIFO has been created
	FifoCreated FifoState = iota
	// FifoAllocated means FIFO has been allocated
	FifoAllocated
)

// String implements fmt.Stringer interface for FifoState
func (fs FifoState) String() string {
	switch fs {
	case FifoCreated:
		return "FIFO_CREATED"
	case FifoAllocated:
		return "FIFO_ALLOCATED"
	default:
		return "FIFO_UNKNOWN_STATE"
	}
}

// FifoOption is FIFO option which can be used to query and set different FIFO properties
// The options starting with RW are both gettable and settable
// The options starting with RO are only gettable
// All settable options except for NC_RWFIFO_HOST_TENSOR_DESCRIPTOR must be set before FIFO is allocated
type FifoOption int

const (
	// RWFifoType configure the fifo type to either of FifoType options
	RWFifoType FifoOption = iota
	// RWFifoConsumerCount is number of consumers of elements before the element is removed
	RWFifoConsumerCount
	// RWFifoDataType configures fifo data type to either of FifoDataType options
	RWFifoDataType
	// RWFifoDontBlock configures to return StatusOutOfMemory instead of blocking
	RWFifoNoBlock
	// ROFifoCapacity allows to query number of maximum elements in the buffer
	ROFifoCapacity
	// ROFifoReadFillLevel allows to query number of tensors in the read buffer
	ROFifoReadFillLevel
	// ROFifoWriteFillLevel allows to query number of tensors in a write buffer
	ROFifoWriteFillLevel
	// ROFifoGraphTensorDescriptor allows to query the tensor descriptor of the FIFO
	ROFifoGraphTensorDesc
	// ROFifoState allows to query FifoState
	ROFifoState
	// ROFifoName allows to query FIFO name
	ROFifoName
	// ROFifoElemDataSize allows to query element data size in bytes
	ROFifoElemDataSize
	// RWFifoHostTensorDesc is tensor descriptor, defaults to none strided channel minor
	RWFifoHostTensorDesc
)

// String implements fmt.Stringer interface
func (fo FifoOption) String() string {
	switch fo {
	case RWFifoType:
		return "RW FIFO type"
	case RWFifoConsumerCount:
		return "RW_FIFO_CONSUMER_COUNT"
	case RWFifoDataType:
		return "RW_FIFO_DATA_TYPE"
	case RWFifoNoBlock:
		return "RW_FIFO_NO_BLOCK"
	case ROFifoCapacity:
		return "RW_FIFO_CAPACITY"
	case ROFifoReadFillLevel:
		return "RO_FIFO_READ_FILL_LEVEL"
	case ROFifoWriteFillLevel:
		return "RO_FIFO_WRITE_FILL_LEVEL"
	case ROFifoGraphTensorDesc:
		return "RO_FIFO_GRAPH_TENSOR_DESCRIPTOR"
	case ROFifoState:
		return "RO_FIFO_STATE"
	case ROFifoName:
		return "RO_FIFO_NAME"
	case ROFifoElemDataSize:
		return "RO_FIFO_ELEM_DATA_SIZE"
	case RWFifoHostTensorDesc:
		return "RW_FIFO_HOST_TENSOR_DESCRIPTOR"
	default:
		return "FIFO_UNKNOWN_OPTION"
	}
}

// Decode decodes raw options data and returns it. The returned data can be asserted into particular type
// It returns error if the data fail to be converted into the option native type
func (fo FifoOption) Decode(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	switch fo {
	case RWFifoType,
		RWFifoConsumerCount,
		RWFifoDataType,
		RWFifoNoBlock,
		ROFifoCapacity,
		ROFifoReadFillLevel,
		ROFifoWriteFillLevel,
		ROFifoElemDataSize,
		ROFifoState:

		var val uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		// this is safe as we expect val to be int
		return uint(val), nil

	case ROFifoName:
		return string(data), nil

	default:
		var val struct {
			BatchSize uint32
			Channels  uint32
			Width     uint32
			Height    uint32
			Size      uint32
			CStride   uint32
			WStride   uint32
			HStride   uint32
			DataType  int32
		}

		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return &TensorDesc{
			BatchSize: uint(val.BatchSize),
			Channels:  uint(val.Channels),
			Width:     uint(val.Width),
			Height:    uint(val.Height),
			Size:      uint(val.Size),
			CStride:   uint(val.CStride),
			WStride:   uint(val.WStride),
			HStride:   uint(val.HStride),
			DataType:  FifoDataType(val.DataType),
		}, nil
	}
}

// FifoOpts specifies FIFO configuration options
type FifoOpts struct {
	// Type is FIFO type
	Type FifoType
	// DataType is FIFO data type
	DataType FifoDataType
	// NumElem is a max number of elements that the FIFO will be able to contain
	NumElem int
}

// Fifo is NCSDK FIFO queue
type Fifo struct {
	name   string
	handle unsafe.Pointer
	device *Device
}

// NewFifo creates new FIFO queue with given name and returns it
// It returns error if it fails to create new queue
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoCreate.html
func NewFifo(name string, t FifoType) (*Fifo, error) {
	var handle unsafe.Pointer

	_name := C.CString(name)
	defer C.free(unsafe.Pointer(_name))

	s := C.ncs_FifoCreate(_name, C.ncFifoType(t), &handle)

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to create new FIFO: %s", Status(s))
	}

	return &Fifo{name: name, handle: handle}, nil
}

// Allocate allocates memory for a FIFO for the specified device based on the number of elements the FIFO will hold and tensorDesc, which describes the expected shape of the FIFO’s elements
// It returns error when it fails to allocate FIFO
//
// More information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoAllocate.html
func (f *Fifo) Allocate(d *Device, td *TensorDesc, numElem uint) error {
	_td := C.struct_ncTensorDescriptor_t{
		n:         C.uint(td.BatchSize),
		c:         C.uint(td.Channels),
		w:         C.uint(td.Width),
		h:         C.uint(td.Height),
		totalSize: C.uint(td.Size),
		cStride:   C.uint(td.CStride),
		wStride:   C.uint(td.WStride),
		hStride:   C.uint(td.HStride),
		dataType:  C.ncFifoDataType(td.DataType),
	}

	s := C.ncs_FifoAllocate(f.handle, d.handle, &_td, C.uint(numElem))

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to allocate FIFO: %s", Status(s))
	}

	return nil
}

// GetOptions queries FIFO options and returns it encoded in a byte slice
// It returns error if it fails to retrieve the options
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoGetOption.html
func (f *Fifo) GetOption(opt FifoOption) ([]byte, error) {
	if opt == RWFifoNoBlock {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	var data unsafe.Pointer
	var dataLen C.uint

	s := C.ncs_FifoGetOption(f.handle, C.int(opt), data, &dataLen)

	if Status(s) == StatusInvalidDataLength {
		return f.GetOptionWithSize(opt, uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option", opt)
}

// GetOptionsWithSize queries FIFO options and returns it encoded in a byte slice of the same size as requested if possible. This function is similar to GetOption(), however as opposed to figuring out the byte size of the queried options it attempts to request the options data by specifying its size explicitly. Because we specify the options data size explicitly this function returns the options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoGetOption.html
func (f *Fifo) GetOptionWithSize(opt FifoOption, size uint) ([]byte, error) {
	if opt == RWFifoNoBlock {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("fifo", f.handle, int(opt), size)
}

// WriteElem writes an element to a FIFO, usually an input tensor for inference along with some metadata
// If it fails to write the element it returns error
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoWriteElem.html
func (f *Fifo) WriteElem(data []byte, metaData interface{}) error {
	dataLen := C.uint(len(data))

	s := C.ncs_FifoWriteElem(f.handle, unsafe.Pointer(&data[0]), &dataLen, unsafe.Pointer(&metaData))

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to write FIFO element: %s", Status(s))
	}

	return nil
}

// ReadElem reads an element from a FIFO, usually the result of an inference as a tensor, along with the associated user-defined data
// If it fails to read the element it returns error
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoReadElem.html
func (f *Fifo) ReadElem() (*Tensor, error) {
	opts, err := f.GetOptionWithSize(ROFifoElemDataSize, C.sizeof_int)
	if err != nil {
		return nil, err
	}

	elemSize, err := ROFifoElemDataSize.Decode(opts)
	if err != nil {
		return nil, err
	}

	var metaData unsafe.Pointer
	size := C.uint(elemSize.(uint))
	data := C.malloc(C.sizeof_char * C.ulong(elemSize.(uint)))

	s := C.ncs_FifoReadElem(f.handle, data, &size, &metaData)

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to read FIFO element: %s", Status(s))
	}

	return &Tensor{
		Data: C.GoBytes(data, C.int(size)),
	}, nil
}

// RemoveElem removes an element from a FIFO
// If it fails to remove the element it returns error
// THIS FUNCTION IS NOT IMPLEMENTED YET
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoRemoveElem.html
func (f *Fifo) RemoveElem() error {
	return fmt.Errorf("%s", StatusUnsupportedFeature)
}

// Destroy destroys NCS FIFO handle and frees associated resources.
// This function must be called for every FIFO handle that was initialized with NewFifo()
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoDestroy.html
func (f *Fifo) Destroy() error {
	s := C.ncs_FifoDestroy(&f.handle)

	if Status(s) != StatusOK {
		return fmt.Errorf("Failed to destroy FIFO: %s", Status(s))
	}

	return nil
}
