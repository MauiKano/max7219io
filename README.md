MAX7219 driver and 8 Digit 7 Segment Display
============================================

This library written in [Go programming language](https://golang.org/) to output a number string to SPI Max 7219 x 8 Seven Segment Display

This branch here is a fork of:

![image](https://raw.github.com/talkkonnect/max7219/master/images/max7219.jpg)

The modification of the fork is the replacement of the backend SDL driver. Here in this fork

periph.ip (periph.io/x/conn/v3/spi)

is used instead of 

spidev (github.com/fulr/spidev)

which is used in the original version.

Compatibility
-------------
Tested on Raspberry PI 3 (model B+)

Golang usage
------------

```go
package main

import (
	"log"
	"github.com/talkkonnect/max7219"
)

func main() {
	mtx := max7219.NewMatrix(1)
	err := mtx.Open(0, 0, 7)
	if err != nil {
		log.Fatal(err)
	}
	defer mtx.Close()

	mtx.Device.SevenSegmentDisplay("1234")
}
```

Installation
------------

```bash
$ go get -u github.com/talkkonnect/max7219
```

Credits
-------

This project is mainly a fork of https://github.com/d2r2/go-max7219

Contact
-------

Please use [Github issue tracker](https://github.com/talkkonnect/max7219/issues) for filing bugs or feature requests.

License
-------

Go-max7219 is licensed under MIT License.


