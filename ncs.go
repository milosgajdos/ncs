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

// StatusCode is the NCSDK API status code as returned by most API calls.
// It usually reports the status of the neural compute stick.
type StatusCode int

const (
	// StatusOK means the API function call worked as expected
	StatusOK StatusCode = -iota
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
func (nc StatusCode) String() string {
	switch nc {
	case StatusOK:
		return "OK"
	case StatusBusy:
		return "Device busy"
	case StatusError:
		return "Unexpected error"
	case StatusOutOfMemory:
		return "Host out of memory"
	case StatusDeviceNotFound:
		return "Device not found"
	case StatusInvalidParameters:
		return "Invalid parameters"
	case StatusTimeout:
		return "Device timeout"
	case StatusCmdNotFound:
		return "Device bootloader not found"
	case StatusNotAllocated:
		return "Unallocated resource"
	case StatusUnauthorized:
		return "Unauthorized operation"
	case StatusUnsupportedGraphFile:
		return "Unsupported graph file"
	case StatusUnsupportedConfigFile:
		return "Unsupported configuration"
	case StatusUnsupportedFeature:
		return "Unsupported feature"
	case StatusMyriadError:
		return "Movidius VPU failure"
	case StatusInvalidDataLength:
		return "Invalid data length when querying options"
	case StatusInvalidHandle:
		return "Invalid handle"
	default:
		return "Unknown"
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

	c := C.ncs_DeviceCreate(C.int(index), &handle)

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to create new device: %s", StatusCode(c))
	}

	return &Device{handle: handle}, nil
}

// Open initializes NCS device and opens device communication channel
// It returns error if it fails to open the device.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceOpen.html
func (d *Device) Open() error {
	c := C.ncs_DeviceOpen(d.handle)

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to open device: %s", StatusCode(c))
	}

	return nil
}

// Close closes the communication channel with NCS device.
// It returns error if it fails to close the communication channel.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceClose.html
func (d *Device) Close() error {
	c := C.ncs_DeviceClose(d.handle)

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to close device: %s", StatusCode(c))
	}

	return nil
}

// Destroy destroys NCS device handle and frees associated resources.
// This function must be called for every device that was initialized with NewDevice().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceDestroy.html
func (d *Device) Destroy() error {
	c := C.ncs_DeviceDestroy(&d.handle)

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to destroy device: %s", StatusCode(c))
	}

	return nil
}

// Graph is NCSDK neural network graph
type Graph struct {
	name    string
	handle  unsafe.Pointer
	device  *Device
	inFifo  *Fifo
	outFifo *Fifo
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

	c := C.ncs_GraphCreate(_name, &handle)

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to create new graph: %s", StatusCode(c))
	}

	return &Graph{name: name, handle: handle}, nil
}

// Allocate allocates a graph on NCS device. This function sends graphData to NCS device. It does not allocate input or output FIFO queues. You have to either allocate them separately or use either AllocateWithFifosDefault() or AllocateWithFifosOpts() functions whcih conveniently create and allocate the FIFO queues.
// It returns error if it fails to allocate the graph on the device
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphAllocate.html
func (g *Graph) Allocate(d *Device, graphData []byte) error {
	c := C.ncs_GraphAllocate(d.handle, g.handle, unsafe.Pointer(&graphData[0]), C.uint(len(graphData)))

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to allocate new graph: %s", StatusCode(c))
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

	c := C.ncs_GraphAllocateWithFifosEx(d.handle,
		g.handle, unsafe.Pointer(&graphData[0]), C.uint(len(graphData)),
		&inHandle, C.ncFifoType(inOpts.Type), C.int(inOpts.NumElem), C.ncFifoDataType(inOpts.DataType),
		&outHandle, C.ncFifoType(outOpts.Type), C.int(outOpts.NumElem), C.ncFifoDataType(outOpts.DataType))

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to allocate graph with FIFOs: %s", StatusCode(c))
	}

	inFifo := &Fifo{handle: inHandle, device: d}
	outFifo := &Fifo{handle: outHandle, device: d}

	g.device = d
	g.inFifo = inFifo
	g.outFifo = outFifo

	return &FifoQueue{
		In:  g.inFifo,
		Out: g.outFifo,
	}, nil
}

// Destroy destroys NCS graph handle and frees associated resources.
// This function must be called for every graph that was initialized with NewGraph().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphDestroy.html
func (g *Graph) Destroy() error {
	c := C.ncs_GraphDestroy(&g.handle)

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to destroy graph: %s", StatusCode(c))
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

// FifoState is state of FIFO
type FifoState int

const (
	// FifoCreated means FIFO has been created
	FifoCreated FifoState = iota
	// FifoAllocated means FIFO has been allocated
	FifoAllocated
)

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
		return "RW_FIFO_Consumer_Count"
	case RWFifoDataType:
		return "RW_FIFO_Data_Type"
	case RWFifoNoBlock:
		return "RW_FIFO_No_Block"
	case ROFifoCapacity:
		return "RW_FIFO_Capacity"
	case ROFifoReadFillLevel:
		return "RO_FIFO_Read_Fill_Level"
	case ROFifoWriteFillLevel:
		return "RO_FIFO_Write_Fill_Level"
	case ROFifoGraphTensorDesc:
		return "RO_FIFO_Graph_Tensor_Descriptor"
	case ROFifoState:
		return "RO_FIFO_State"
	case ROFifoName:
		return "RO_FIFO_Name"
	case ROFifoElemDataSize:
		return "RO_FIFO_Elem_Data_Size"
	case RWFifoHostTensorDesc:
		return "RW_FIFO_Host_Tensor_Descriptor"
	default:
		return "Unknown"
	}
}

// Decode decodes raw options data and returns it. The returned data can be asserted into particular type
// It returns error if the data fail to be converted into the option native type
func (fo FifoOption) Decode(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	switch fo {
	case RWFifoType, RWFifoConsumerCount, RWFifoDataType, RWFifoNoBlock, ROFifoCapacity,
		ROFifoReadFillLevel, ROFifoWriteFillLevel, ROFifoElemDataSize, ROFifoState:
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

	return nil, nil
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

	c := C.ncs_FifoCreate(_name, C.ncFifoType(t), &handle)

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to create new FIFO: %s", StatusCode(c))
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

	c := C.ncs_FifoAllocate(f.handle, d.handle, &_td, C.uint(numElem))

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to allocate FIFO: %s", StatusCode(c))
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
		return nil, fmt.Errorf("Not implemented")
	}

	optsData := C.OptionsData{}

	c := C.ncs_FifoGetOption(f.handle, C.int(opt), &optsData)

	if StatusCode(c) == StatusInvalidDataLength {
		// allocate the data with correct size and try again
		optsData.data = C.malloc(C.sizeof_char * C.ulong(optsData.length))
		defer C.free(unsafe.Pointer(optsData.data))

		c = C.ncs_FifoGetOption(f.handle, C.int(opt), &optsData)

		if StatusCode(c) != StatusOK {
			return nil, fmt.Errorf("Failed to get FIFO options: %s", StatusCode(c))
		}
	}

	return C.GoBytes(unsafe.Pointer(optsData.data), C.int(optsData.length)), nil
}

// GetOptionsWithSize queries FIFO options and returns it encoded in a byte slice of the same size as requested if possible. This function is similar to GetOption(), however as opposed to figuring out the byte size of the queried options it attempts to request the options data by specifying its size explicitly. Because we specify the options data size explicitly this function returns the options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoGetOption.html
func (f *Fifo) GetOptionWithSize(opt FifoOption, size uint) ([]byte, error) {
	if opt == RWFifoNoBlock {
		return nil, fmt.Errorf("Not implemented")
	}

	optsData := C.OptionsData{}
	optsData.length = C.uint(size)

	c := C.ncs_FifoGetOption(f.handle, C.int(opt), &optsData)

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to get FIFO options: %s", StatusCode(c))
	}

	return C.GoBytes(unsafe.Pointer(optsData.data), C.int(optsData.length)), nil
}

// WriteElem writes an element to a FIFO, usually an input tensor for inference along with some metadata
// If it fails to write the element it returns error
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoWriteElem.html
func (f *Fifo) WriteElem(data []byte, metaData interface{}) error {
	dataLen := C.uint(len(data))
	c := C.ncs_FifoWriteElem(f.handle, unsafe.Pointer(&data[0]), &dataLen, unsafe.Pointer(&metaData))

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to write FIFO element: %s", StatusCode(c))
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
	//opts, err := f.GetOption(ROFifoElemDataSize)
	if err != nil {
		return nil, err
	}

	elemSize, err := ROFifoElemDataSize.Decode(opts)
	if err != nil {
		return nil, err
	}

	var data unsafe.Pointer
	var metaData unsafe.Pointer
	size := C.uint(elemSize.(uint))

	c := C.ncs_FifoReadElem(f.handle, data, &size, &metaData)

	if StatusCode(c) != StatusOK {
		return nil, fmt.Errorf("Failed to read FIFO element: %s", StatusCode(c))
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
	c := C.ncs_FifoDestroy(&f.handle)

	if StatusCode(c) != StatusOK {
		return fmt.Errorf("Failed to destroy FIFO: %s", StatusCode(c))
	}

	return nil
}
