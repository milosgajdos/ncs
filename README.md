# ncs
Neural Compute Stick V2 API Go binding

So far it only contains Device API implementation.

**WARNING, NCSDK API V2 IS BADLY BROKEN AT THE MOMENT**

# Quick start

On MacOS, clone macos branch:

```shell
$ git clone -b macos https://github.com/milosgajdos83/ncsdk.git
```

Build ncsdk API libraries:

```shell
$ cd api/src && sudo make basicinstall pythoninstall
```

Test NCSDK example:

```shell
$ cd ../../examples/apps/hello_ncs_cpp/ && make run
```

# Example Go program

The example below shows how to created and destroy the basic types the NCSDK API 2.0 provides

```go
package main

import (
	"fmt"
	"os"

	"github.com/milosgajdos83/ncs"
)

func ExitOnErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Creating NCS device handle")
	dev, err := ncs.NewDevice(0)
	ExitOnErr(err)
	fmt.Println("NCS device handle created")

	fmt.Println("Opening NCS device")
	err = dev.Open()
	ExitOnErr(err)
	fmt.Println("NCS device opened")

	fmt.Println("Creating NCS graph handle")
	graph, err := ncs.NewGraph("TestGraph")
	ExitOnErr(err)
	fmt.Println("NCS graph created")

	fmt.Println("Creating NCS FIFO handle")
	fifo, err := ncs.NewFifo("TestFIFO", ncs.FifoHostRO)
	ExitOnErr(err)
	fmt.Println("NCS FIFO handle created")

	fmt.Println("Destroying NCS FIFO")
	err = fifo.Destroy()
	ExitOnErr(err)
	fmt.Println("NCS FIFO destroyed")

	fmt.Println("Destroyig NCS graph")
	err = graph.Destroy()
	ExitOnErr(err)
	fmt.Println("NCS graph destroyed")

	fmt.Println("Closing NCS device")
	err = dev.Close()
	ExitOnErr(err)
	fmt.Println("NCS device closed")

	fmt.Println("Destroyig NCS device")
	err = dev.Destroy()
	ExitOnErr(err)
	fmt.Println("NCS device destroyed")
}
```

If your Movidius NCS device is plugged in you should see the following output when running the program above:

```console
Creating NCS device handle
NCS device handle created
Opening NCS device
NCS device opened
Creating NCS graph handle
NCS graph created
Creating NCS FIFO handle
NCS FIFO handle created
Destroying NCS FIFO
NCS FIFO destroyed
Destroyig NCS graph
NCS graph destroyed
Closing NCS device
NCS device closed
Destroyig NCS device
NCS device destroyed
```
