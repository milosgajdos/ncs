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

// FifoQueue is a FIFO queue used for NCS inference.
type FifoQueue struct {
	// In is an inbound queue
	In *Fifo
	// Out is an outbound queue
	Out *Fifo
}

// FifoType defines FIFO access types.
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

// FifoDataType defines possible data types for FIFOs.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoDataType_t.html
type FifoDataType int

const (
	// FifoFP16 data is in half precision (16 bit) floating point format (FP16).
	FifoFP16 FifoDataType = iota
	// FifoFP32 data is in full precision (32 bit) floating point format (FP32).
	FifoFP32
)

// String implements fmt.Stringer interface
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

// FifoState represents FIFO state
type FifoState int

const (
	// FifoCreated means FIFO has been created.
	FifoCreated FifoState = iota
	// FifoAllocated means FIFO has been allocated.
	FifoAllocated
)

// String implements fmt.Stringer interface
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

// FifoOption is FIFO option which can be used to query and set various FIFO properties.
// The options starting with RW are both gettable and settable.
// The options starting with RO are only gettable.
// All settable options except for RWFifoHostTensorDesc must be set before FIFO is allocated.
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

// fifoOptSize is a map which maps FIFO options to its native sizes
var fifoOptSize = map[Option]uint{
	RWFifoType:            C.sizeof_int,
	RWFifoConsumerCount:   C.sizeof_int,
	RWFifoDataType:        C.sizeof_int,
	RWFifoNoBlock:         C.sizeof_int,
	ROFifoCapacity:        C.sizeof_int,
	ROFifoReadFillLevel:   C.sizeof_int,
	ROFifoWriteFillLevel:  C.sizeof_int,
	ROFifoGraphTensorDesc: C.sizeof_struct_ncTensorDescriptor_t,
	ROFifoState:           C.sizeof_int,
	ROFifoName:            C.sizeof_char,
	ROFifoElemDataSize:    C.sizeof_int,
	RWFifoHostTensorDesc:  C.sizeof_struct_ncTensorDescriptor_t,
}

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

// Value returns option value as integer
func (fo FifoOption) Value() int {
	return int(fo)
}

// Decode decodes options data encoded in raw bytes and returns it in its native type.
// The returned data can be asserted into its native type.
// If the data contains more than one element you need to specify the number of expected elements via count.
// It returns error if the data fails to be decoded into the option native type.
func (fo FifoOption) Decode(data []byte, count int) (interface{}, error) {
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

		// this is safe type cast as we know val is positive integer
		return uint(val), nil

	case ROFifoName:
		return string(data), nil

	case RWFifoHostTensorDesc:
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

	default:
		return nil, fmt.Errorf("Unable to decode FIFO option data: %s", fo)
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
		return f.GetOptionWithByteSize(opt, fifoOptSize[opt]*uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option: %s", opt, Status(s))
}

// GetOptionsWithSize queries NCS fifo options and returns it encoded in a byte slice of size elements.
// This function is similar to GetOption(), however as opposed to GetOption() which first queries the NCS device for the size of the requested options, it attempts to request the options data by specifying its size in raw bytes explicitly, hence it returns the queried options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncFifoGetOption.html
func (f *Fifo) GetOptionWithByteSize(opt FifoOption, size uint) ([]byte, error) {
	if opt == RWFifoNoBlock {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("fifo", f.handle, opt, size)
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
	opts, err := f.GetOptionWithByteSize(ROFifoElemDataSize, C.sizeof_int)
	if err != nil {
		return nil, err
	}

	elemSize, err := ROFifoElemDataSize.Decode(opts, 1)
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
