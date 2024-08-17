/*
 * talkkonnect headless mumble client/gateway with lcd screen and channel control
 * Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * talkkonnect is the based on talkiepi and barnard by Daniel Chote and Tim Cooper
 *
 * The Initial Developer of the Original Code is
 * Suvir Kumar <suvir@talkkonnect.com>
 * Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
 *
 * Original Code was copied from https://github.com/d2r2/go-max7219
 *
 * Original Code with MIT License
 *
 * Copyright (c) 2015 Richard Hull
 * Copyright (c) 2015 Denis Dyakov
 * Contributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 * output_sevensegment.go modified to work with talkkonnect
 */

package max7219

import (
	"fmt"
	"log"
	"spidev"
)

type Max7219Reg byte

const (
	MAX7219_REG_NOOP   Max7219Reg = 0
	MAX7219_REG_DIGIT0            = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
	MAX7219_REG_LASTDIGIT   = MAX7219_REG_DIGIT7
)

const MAX7219_DIGIT_COUNT = MAX7219_REG_LASTDIGIT -
	MAX7219_REG_DIGIT0 + 1

type Device struct {
	cascaded int
	buffer   []byte
	spi      *spidev.SPIDevice
}

func NewDevice(cascaded int) *Device {
	buf := make([]byte, MAX7219_DIGIT_COUNT*cascaded)
	this := &Device{cascaded: cascaded, buffer: buf}
	return this
}

func (this *Device) GetCascadeCount() int {
	return this.cascaded
}

func (this *Device) GetLedLineCount() int {
	return MAX7219_DIGIT_COUNT
}

func (this *Device) Open(spibus int, spidevice int, brightness byte) error {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", spibus, spidevice)
        fmt.Println("max7219 open device", devstr)
	spi, err := spidev.NewSPIDevice(devstr)
	if err != nil {
		return err
	}
	this.spi = spi
	// Initialize Max7219 led driver.
	this.Command(MAX7219_REG_SCANLIMIT, 7)   // show all 8 digits
	this.Command(MAX7219_REG_DECODEMODE, 0)  // use matrix (not digits)
	this.Command(MAX7219_REG_DISPLAYTEST, 0) // no display test
	this.Command(MAX7219_REG_SHUTDOWN, 1)    // not shutdown mode
	this.Brightness(brightness)
	this.ClearAll(true)
	return nil
}

func (this *Device) Close() {
	this.spi.Close()
}

func (this *Device) Brightness(intensity byte) error {
	return this.Command(MAX7219_REG_INTENSITY, intensity)
}

func (this *Device) Command(reg Max7219Reg, value byte) error {
	buf := []byte{byte(reg), value}
	for i := 0; i < this.cascaded; i++ {
		_, err := this.spi.Xfer(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) sendBufferLine(position int) error {
	reg := MAX7219_REG_DIGIT0 + position
	//fmt.Printf("Register: %#x\n", reg)
	buf := make([]byte, this.cascaded*2)
	for i := 0; i < this.cascaded; i++ {
		b := this.buffer[i*MAX7219_DIGIT_COUNT+position]
		//fmt.Printf("Buffer value: %#x\n", b)
		buf[i*2] = byte(reg)
		buf[i*2+1] = b
	}
	//log.Printf("debug: Send to bus: %v\n", buf)
	_, err := this.spi.Xfer(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this *Device) SetBufferLine(cascadeId int,
	position int, value byte, redraw bool) error {
	this.buffer[cascadeId*MAX7219_DIGIT_COUNT+position] = value
	if redraw {
		err := this.sendBufferLine(position)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) Flush() error {
	for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
		err := this.sendBufferLine(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) Clear(cascadeId int, redraw bool) error {
	if cascadeId >= 0 {
		for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
			this.buffer[cascadeId*MAX7219_DIGIT_COUNT+i] = 0
		}
	} else {
		for i := 0; i < this.cascaded*MAX7219_DIGIT_COUNT; i++ {
			this.buffer[i] = 0
		}
	}
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ClearAll(redraw bool) error {
	for i := 0; i < this.cascaded; i++ {
		err := this.Clear(i, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ScrollLeft(redraw bool) error {
	this.buffer = append(this.buffer[1:], 0)
	log.Printf("debug: Buffer: %v\n", this.buffer)
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) ScrollRight(redraw bool) error {
	this.buffer = append([]byte{0}, this.buffer[:len(this.buffer)-1]...)
	if redraw {
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) SevenSegmentDisplay(toDisplay string) error {
	if len(toDisplay) == 0 || len(toDisplay) > 8 {
		return nil
	}

	this.ClearAll(true)

	for i, digit := range toDisplay {
		switch string(digit) {
		case "A", "a":
			this.Command(Max7219Reg(len(toDisplay)-i),119) //A
		case "B", "b":
			this.Command(Max7219Reg(len(toDisplay)-i),31)  //b
		case "C":
			this.Command(Max7219Reg(len(toDisplay)-i),78)  //C
		case "c":
			this.Command(Max7219Reg(len(toDisplay)-i),13)  //c
		case "D", "d":
			this.Command(Max7219Reg(len(toDisplay)-i),61)  //d
		case "E", "e":
			this.Command(Max7219Reg(len(toDisplay)-i),79)  //E
		case "F", "f":
			this.Command(Max7219Reg(len(toDisplay)-i),71)  //F
		case "G", "g":
			this.Command(Max7219Reg(len(toDisplay)-i),94)  //G
		case "H":
			this.Command(Max7219Reg(len(toDisplay)-i),55)  //H
		case "h":
			this.Command(Max7219Reg(len(toDisplay)-i),23)  //h
		case "I", "i":
			this.Command(Max7219Reg(len(toDisplay)-i),6)   //i
		case "L", "l":
			this.Command(Max7219Reg(len(toDisplay)-i),14)  //L
		case "N", "n":
			this.Command(Max7219Reg(len(toDisplay)-i),21)  //n
		case "O":
			this.Command(Max7219Reg(len(toDisplay)-i),99)  //O
		case "o":
			this.Command(Max7219Reg(len(toDisplay)-i),29)  //o
		case "P", "p":
			this.Command(Max7219Reg(len(toDisplay)-i),103) //p 
		case "R", "r":
			this.Command(Max7219Reg(len(toDisplay)-i),5)   //r
		case "S", "s":
			this.Command(Max7219Reg(len(toDisplay)-i),91)  //S
		case "T", "t":
			this.Command(Max7219Reg(len(toDisplay)-i),15)  //t
		case "U":
			this.Command(Max7219Reg(len(toDisplay)-i),62)  //U
		case "u":
			this.Command(Max7219Reg(len(toDisplay)-i),28)  //u
		case "Y", "y":
			this.Command(Max7219Reg(len(toDisplay)-i),59)   //Y
		case "0":
			this.Command(Max7219Reg(len(toDisplay)-i), 126) //0
		case "1":
			this.Command(Max7219Reg(len(toDisplay)-i), 48)  //1 
		case "2":
			this.Command(Max7219Reg(len(toDisplay)-i), 109) //2
		case "3":
			this.Command(Max7219Reg(len(toDisplay)-i), 121) //3
		case "4":
			this.Command(Max7219Reg(len(toDisplay)-i), 51)  //4
		case "5":
			this.Command(Max7219Reg(len(toDisplay)-i), 91)  //5
		case "6":
			this.Command(Max7219Reg(len(toDisplay)-i), 95)  //6
		case "7":
			this.Command(Max7219Reg(len(toDisplay)-i), 112) //7
		case "8":
			this.Command(Max7219Reg(len(toDisplay)-i), 127) //8
		case "9":
			this.Command(Max7219Reg(len(toDisplay)-i), 115) //9
		case ".":
			this.Command(Max7219Reg(len(toDisplay)-i),128)  //.
		case "-":
			this.Command(Max7219Reg(len(toDisplay)-i),129)  //-
		case "=":
			this.Command(Max7219Reg(len(toDisplay)-i),9)    //= lower
		case "J","j":
			this.Command(Max7219Reg(len(toDisplay)-i),64)   //= upper
		case "K","k":
			this.Command(Max7219Reg(len(toDisplay)-i),72)   //= split
		case "M","m":
			this.Command(Max7219Reg(len(toDisplay)-i),73)   //==

		}
	}
	return nil
}
