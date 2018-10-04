# ncs

[![GoDoc](https://godoc.org/github.com/milosgajdos83/ncs?status.svg)](https://godoc.org/github.com/milosgajdos83/ncs)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Neural Compute Stick V2.0 API Go binding

**NCSDK API V2 IS PARTIALLY BROKEN on macOS AT THE MOMENT -- EVERYTHING WORKS FINE ON LINUX**

The code in this repository has been tested on the following Linux OS:

```
Distributor ID:	Ubuntu
Description:	Ubuntu 16.04.5 LTS
Release:	16.04
Codename:	xenial

Linux ubuntu-xenial 4.4.0-134-generic #160-Ubuntu SMP Wed Aug 15 14:58:00 UTC 2018 x86_64 x86_64 x86_64 GNU/Linux
```

The Movidius NCSDK API coverage provided in this repo should give you all the tools to use Movidius NCS to perform Neural Network inference.

# Quick start

On MacOS, clone `macos-V2` branch:

```shell
$ git clone -b macos-V2 https://github.com/milosgajdos83/ncsdk.git
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

The example below shows how to create and destroy the basic resource types the NCSDK API 2.0 provides. For more complex examples please see [examples](./examples)

```go
package main

import (
	"log"

	"github.com/milosgajdos83/ncs"
)

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

	log.Printf("Attempting to create NCS FIFO handle")
	fifo, err := ncs.NewFifo("TestFIFO", ncs.FifoHostRO)
	defer fifo.Destroy()
	if err != nil {
		return
	}
	log.Printf("NCS FIFO handle successfully created")
}
```

If your Movidius NCS device is plugged in you should see the following output when running the program above:

```console
2018/08/27 00:43:00 Attempting to create NCS device handle
2018/08/27 00:43:00 NCS device handle successfully created
2018/08/27 00:43:00 Attempting to open NCS device
2018/08/27 00:43:03 NCS device successfully opened
2018/08/27 00:43:03 Attempting to create NCS graph handle
2018/08/27 00:43:03 NCS graph handle successfully created
2018/08/27 00:43:03 Attempting to create NCS FIFO handle
2018/08/27 00:43:03 NCS FIFO handle successfully created
```
