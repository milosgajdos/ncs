package ncs

// #cgo LDFLAGS: -lmvnc
/*
#include <ncs.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	// MaxNameSize is the maximum length of device or graph name size
	MaxNameSize = 28
	// ThermalBufferSize is the size of the temperature buffer as returned when querying device
	ThermalBufferSize = 100
	// DebugBufferSize is the size of the debug information buffer as returned by API
	DebugBufferSize = 120
	// VersionMaxSize is the max length of various version options (HW, firmewre etc.) as returned by API
	VersionMaxSize = 4
)

// Status is the NCSDK API status code as returned by most API calls.
// It usually reports the status of the Neural Compute Stick.
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
	// StatusNotAllocated means the graph or fifo has not been allocated..
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
	// StatusInvalidDataLength means an invalid data length has been passed when getting or setting an option.
	StatusInvalidDataLength
	// StatusInvalidHandle means an invalid handle has been passed to a function.
	StatusInvalidHandle
)

// String method implements fmt.Stringer interface
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

// Option is NCS option
type Option interface {
	// Value returns Option value as its integer code
	Value() int
	// Decode decodes raw byte slice option data as returned from NCS
	Decode([]byte, int) (interface{}, error)
}

// TensorDesc describes NCS graph inputs and outputs
type TensorDesc struct {
	// BatchSize contains number of elements.
	BatchSize uint
	// Channels contains number of channels (when dealing with digital images).
	Channels uint
	// Width is data width (i.e. number of matrix columns).
	Width uint
	// Height is data height (i.e. number of matrix rows).
	Height uint
	// Size is the total data size in the tensor.
	Size uint
	// CStride is channel stride (Stride in the channels' dimension).
	CStride uint
	// WStride is width stride (Stride in the horizontal dimension).
	WStride uint
	// HStride is height stride (Stride in the vertical dimension).
	HStride uint
	// DataType is data type of the tensor.
	DataType FifoDataType
}

// Tensor is graph tensor as returned from NCS
type Tensor struct {
	// Data contains raw tensor data
	Data []byte
	// MetaData contains tensor metadata
	MetaData interface{}
}

// getOption is a function which unifies querying of various NCS resource options
func getOption(resource string, handle unsafe.Pointer, option Option, size uint) ([]byte, error) {
	// allocate buffer for options data
	data := C.malloc(C.sizeof_char * C.ulong(size))
	defer C.free(unsafe.Pointer(data))
	dataLen := C.uint(size)

	// NCCS API status code
	var s C.int

	switch resource {
	case "device":
		s = C.ncs_DeviceGetOption(handle, C.int(option.Value()), data, &dataLen)
	case "graph":
		s = C.ncs_GraphGetOption(handle, C.int(option.Value()), data, &dataLen)
	case "fifo":
		s = C.ncs_FifoGetOption(handle, C.int(option.Value()), data, &dataLen)
	default:
		return nil, fmt.Errorf("Unknown resource: %s", resource)
	}

	if Status(s) != StatusOK {
		return nil, fmt.Errorf("Failed to get %s option: %s", resource, Status(s))
	}

	return C.GoBytes(unsafe.Pointer(data), C.int(size)), nil
}
