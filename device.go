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

// DeviceHWVersion defines neural compute device hardware version
type DeviceHWVersion int

const (
	MA2450 DeviceHWVersion = iota
	MA2480
)

// String implements fmt.Stringer interface
func (hw DeviceHWVersion) String() string {
	switch hw {
	case MA2450:
		return "MA2450"
	case MA2480:
		return "MA2480"
	default:
		return "UNKNOWN_DEVICE_VERSION"
	}
}

// DeviceThermalThrottle defines thermal throttle level
type DeviceThermalThrottle int

const (
	// NoThrottle means no limit reached
	NoThrottle DeviceThermalThrottle = iota
	// LowerGuard means means lower guard temperature threshold of chip sensor has been reached
	// Short throttling time is in action between inferences to protect the device
	LowerGuard
	// UpperGuards means upper guard temperature of chip sensor has been reached
	// Long throttling time is in action between inferences to protect the device
	UpperGuard
)

// String implements fmt.Stringer interface
func (dt DeviceThermalThrottle) String() string {
	switch dt {
	case NoThrottle:
		return "NO_THERMAL_THROTTLE"
	case LowerGuard:
		return "LOWER_GUARD_THERMAL_THROTTLE"
	case UpperGuard:
		return "UPPER_GUARD_THERMAL_THROTTLE"
	default:
		return "UNKNWON_THERMAL_THROTTLE"
	}
}

// DeviceOption defines NCS device options.
// The options starting with RW are both gettable and settable.
// The options starting with RO are only gettable.
type DeviceOption int

const (
	// RODeviceThermalStats queries device temperatures in degrees Celsius.
	// This option returns []float64 array of max temperatures for the last ThermalBufferSize seconds.
	RODeviceThermalStats DeviceOption = (2000 + iota)
	// RODeviceThermalThrottling queries temperature throttling level.
	RODeviceThermalThrottle
	// RODeviceState queries the state of the device.
	RODeviceState
	// RODeviceMemoryUsed queries current memory in use on the device in bytes.
	RODeviceMemoryUsed
	// RODeviceMemorySize queries total memory available on the device in bytes.
	RODeviceMemorySize
	// RODeviceMaxFifoCount queries maximum number of FIFOs that can be allocated for the device.
	RODeviceMaxFifoCount
	// RODeviceAllocatedFifoCount queries number of FIFOs currently allocated for the device.
	RODeviceAllocatedFifoCount
	// RODeviceMaxMaxGraphCount queries the maximum number of graphs that can be allocated for the device.
	RODeviceMaxGraphCount
	// RODeviceAllocatedGraphCount queries the number of graphs currently allocated for the device.
	RODeviceAllocatedGraphCount
	// RODeviceClassLimit queries the highest device option class supported.
	RODeviceClassLimit
	// RODeviceFirmwareVersion queries the version of the firmware currently running on the device.
	RODeviceFirmwareVersion
	// RODeviceDebugInfo queries more detailed info when the result of API call is StatusMyriadError.
	RODeviceDebugInfo
	// RODeviceMVTensorVersion queries the version of the mvtensor library that was linked with the API.
	RODeviceMVTensorVersion
	// RODeviceName queries the internal name of the device.
	RODeviceName
	// RODeviceMaxExecutors is reserved for future use.
	RODeviceMaxExecutors
	// RODeviceHWVersion queries the hardware version of the device.
	RODeviceHWVersion
)

// deviceOptSize is a map which maps device options to its native sizes
var deviceOptSize = map[Option]uint{
	RODeviceThermalStats:        C.sizeof_float,
	RODeviceThermalThrottle:     C.sizeof_int,
	RODeviceState:               C.sizeof_int,
	RODeviceMemoryUsed:          C.sizeof_int,
	RODeviceMemorySize:          C.sizeof_int,
	RODeviceMaxFifoCount:        C.sizeof_int,
	RODeviceAllocatedFifoCount:  C.sizeof_int,
	RODeviceMaxGraphCount:       C.sizeof_int,
	RODeviceAllocatedGraphCount: C.sizeof_int,
	RODeviceClassLimit:          C.sizeof_int,
	RODeviceFirmwareVersion:     C.sizeof_uint,
	RODeviceDebugInfo:           C.sizeof_char,
	RODeviceMVTensorVersion:     C.sizeof_uint,
	RODeviceName:                C.sizeof_char,
	RODeviceMaxExecutors:        C.sizeof_int,
	RODeviceHWVersion:           C.sizeof_int,
}

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

// Value returns option value as integer
func (do DeviceOption) Value() int {
	return int(do)
}

// Decode decodes raw options data and returns it. The returned data can be asserted into its native type.
// It returns error if the data fails to be decoded into the option native type.
func (do DeviceOption) Decode(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)

	switch do {
	case RODeviceThermalThrottle,
		RODeviceState,
		RODeviceMemoryUsed,
		RODeviceMemorySize,
		RODeviceMaxFifoCount,
		RODeviceAllocatedFifoCount,
		RODeviceMaxGraphCount,
		RODeviceAllocatedGraphCount,
		RODeviceClassLimit,
		RODeviceMaxExecutors,
		RODeviceHWVersion:

		var val uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		// this is safe type cast as we know val is positive integer
		return uint(val), nil

	case RODeviceThermalStats:

		var val [ThermalBufferSize]float32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return val[:], nil

	case RODeviceFirmwareVersion:

		var val [VersionMaxSize]uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return val[:], nil

	case RODeviceMVTensorVersion:

		var val [2]uint32
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}

		return val[:], nil

	case RODeviceDebugInfo,
		RODeviceName:

		return string(data), nil

	default:
		return nil, fmt.Errorf("Unable to decode device option data: %s", do)
	}
}

// DeviceState represents NCS device state
type DeviceState int

const (
	// DeviceCreated means NCS device handle has been created.
	DeviceCreated DeviceState = iota
	// DeviceOpened means NCS device handle has been opened.
	DeviceOpened
	// DeviceClosed means NCS device handle has been closed.
	DeviceClosed
)

// String implements fmt.Stringer interface
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

// NewDevice creates new NCS device handle and returns it.
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

// Open initializes NCS device and opens device communication channel.
// It returns error if it fails to open or initialize the communication channel with the device.
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

// GetOption queries the value of an option for the device and returns it encoded in a byte slice.
// It returns error if it fails to retrieve the option value.
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
		return d.GetOptionWithByteSize(opt, deviceOptSize[opt]*uint(dataLen))
	}

	return nil, fmt.Errorf("Failed to read %s option: %s", opt, Status(s))
}

// GetOptionsWithSize queries NCS device options and returns it encoded in a byte slice of size elements.
// This function is similar to GetOption(), however as opposed to GetOption() which first queries the NCS device for the size of the requested options, it attempts to request the options data by specifying its size in raw bytes explicitly, hence it returns the queried options data faster.
// It returns error if it fails to retrieve the options or if the requested size of the options is invalid.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceGetOption.html
func (d *Device) GetOptionWithByteSize(opt DeviceOption, size uint) ([]byte, error) {
	if opt == RODeviceMaxExecutors || opt == RODeviceDebugInfo {
		return nil, fmt.Errorf("Option %s not implemented", opt)
	}

	return getOption("device", d.handle, opt, size)
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
