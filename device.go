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

package max7219io

import (
	//	"fmt"

	"fmt"
	"log"

	//	"time"
	//        "spidev"
	//	"github.com/fulr/spidev"
	//	"periph.io/x/conn/v3/driver/driverreg"
	//	"periph.io/x/conn/v3/physic"
	//"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/spi"
	//"periph.io/x/host/v3/rpi"
	// "periph.io/x/conn/v3/spi/spireg"
)

type Max7219Reg byte

const (
	MAX7219_REG_NOOP        = 0x00
	MAX7219_REG_DIGIT0      = 0x01
	MAX7219_REG_DIGIT1      = 0x02
	MAX7219_REG_DIGIT2      = 0x03
	MAX7219_REG_DIGIT3      = 0x04
	MAX7219_REG_DIGIT4      = 0x05
	MAX7219_REG_DIGIT5      = 0x06
	MAX7219_REG_DIGIT6      = 0x07
	MAX7219_REG_DIGIT7      = 0x08
	MAX7219_REG_DECODEMODE  = 0x09
	MAX7219_REG_INTENSITY   = 0x0A
	MAX7219_REG_SCANLIMIT   = 0x0B
	MAX7219_REG_SHUTDOWN    = 0x0C
	MAX7219_REG_DISPLAYTEST = 0x0F
	MAX7219_REG_LASTDIGIT   = MAX7219_REG_DIGIT7

	OP_ON  = 0x01
	OP_OFF = 0x00
	MAX72b = 0x00
	MAX72d = 0x01
)

const MAX7219_DIGIT_COUNT = MAX7219_REG_LASTDIGIT - MAX7219_REG_DIGIT0 + 1

type Device struct {
	cascaded int
	buffer   []byte
	conn     spi.Conn
}

func NewDevice(cascaded int, conn spi.Conn) *Device {
	buf := make([]byte, MAX7219_DIGIT_COUNT*cascaded)
	this := &Device{cascaded: cascaded,
		buffer: buf,
		conn:   conn}
	return this
}

func (this *Device) GetCascadeCount() int {
	return this.cascaded
}

func (this *Device) GetLedLineCount() int {
	return MAX7219_DIGIT_COUNT
}

func (this *Device) Open(brightness byte) error {
	this.Command(MAX7219_REG_SCANLIMIT, 7)   // show all 8 digits
	this.Command(MAX7219_REG_DECODEMODE, 0)  // use individual segments
	this.Command(MAX7219_REG_DISPLAYTEST, 0) // no display test
	//this.Command(MAX7219_REG_SHUTDOWN, OP_OFF) // not shutdown mode
	this.Command(MAX7219_REG_SHUTDOWN, OP_ON) // not shutdown mode
	this.Brightness(brightness)
	this.ClearAll(true)
	//this.Command(MAX7219_REG_NOOP, 0)
	return nil
}

func (this *Device) Close() {
	// this.con.Close()
	// this.conn.Tx
}

// Write sends data to the device via SPI
func (this *Device) Write(data []byte) error {
	//fmt.Printf("Buffer size: %d bytes\n", len(data))
	//cs := rpi.P1_24 // GPIO8 for chip select
	//dout := rpi.P1_19 // GPI10 for dout
	//clk := rpi.P1_23  // GPI11 for CLK
	//clk.Out(gpio.Low)
	//dout.Out(gpio.Low)
	//cs.Out(gpio.Low)
	err := this.conn.Tx(data, nil)
	//cs.Out(gpio.High)
	//clk.Out(gpio.Low)
	//dout.Out(gpio.Low)
	//return this.conn.Tx(data, nil) // Send data to the device
	return err
}

func (this *Device) Brightness(intensity byte) error {
	return this.Command(MAX7219_REG_INTENSITY, intensity)
}

func (this *Device) SetRegisters() {
	this.Command(MAX7219_REG_SCANLIMIT, 7)    // show all 8 digits
	this.Command(MAX7219_REG_DECODEMODE, 0)   // use individual segments
	this.Command(MAX7219_REG_SHUTDOWN, OP_ON) // not shutdown mode
	this.Command(MAX7219_REG_NOOP, 0)
}

func (this *Device) NOP7219() error {
	return this.Command(MAX7219_REG_NOOP, 0)
}

func (this *Device) Command(reg Max7219Reg, value byte) error {
	buf := []byte{byte(reg), value}
	for i := 0; i < this.cascaded; i++ {
		err := this.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) sendBufferLine0() error {
	reg := byte(MAX7219_REG_DIGIT0)
	buf := make([]byte, this.cascaded*6)
	for i := 0; i < this.cascaded; i++ {
		b := this.buffer[i*MAX7219_DIGIT_COUNT]
		//fmt.Printf("Buffer value[%d],%d: %#x\n", i,position,b)
		buf[i*2+4] = byte(MAX7219_REG_DISPLAYTEST)
		buf[i*2+5] = byte(0)
		buf[i*2+2] = byte(MAX7219_REG_SHUTDOWN)
		buf[i*2+3] = byte(OP_ON)

		buf[i*2] = byte(reg)
		buf[i*2+1] = byte(b)
	}
	//log.Printf("debug: Send to bus: %v\n", buf)
	err := this.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this *Device) sendBufferLine(position int) error {
	reg := byte(MAX7219_REG_DIGIT0 + position)
	//buf := make([]byte, this.cascaded*6)
	buf := make([]byte, this.cascaded*2)
	for i := 0; i < this.cascaded; i++ {
		b := this.buffer[i*MAX7219_DIGIT_COUNT+position]
		//fmt.Printf("Buffer value[%d],%d: %#x\n", i,position,b)
		/*
			buf[i*2] = byte(MAX7219_REG_DISPLAYTEST)
			buf[i*2+1] = byte(0)
			buf[i*2+2] = byte(MAX7219_REG_SHUTDOWN)
			buf[i*2+3] = byte(OP_ON)

			buf[i*2+4] = byte(reg)
			buf[i*2+5] = byte(b)
		*/
		buf[i*2] = byte(reg)
		buf[i*2+1] = byte(b)

	}
	//log.Printf("debug: Send to bus: %v\n", buf)
	err := this.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// this.buffer[cascadeId*MAX7219_DIGIT_COUNT+position] = value
func (this *Device) sendBufferFull() error {
	reg := byte(MAX7219_REG_DIGIT0)
	buf := make([]byte, MAX7219_DIGIT_COUNT*this.cascaded*2)
	for j := 0; j < MAX7219_DIGIT_COUNT; j++ {
		reg = byte(MAX7219_REG_DIGIT0 + j)
		for i := 0; i < this.cascaded; i++ {
			b := this.buffer[i*MAX7219_DIGIT_COUNT+j]
			fmt.Printf("Buffer value[%d, lin=%d],%d: %#x, %#x\n", i, j, (2*j*this.cascaded)+(i*2), b, reg)
			buf[(2*j*this.cascaded)+(i*2)] = byte(reg)
			buf[(2*j*this.cascaded)+((i*2)+1)] = byte(b)
		}
	}
	err := this.Write(buf)
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

/*
func (this *Device) Flush() error {
	err := this.sendBufferFull()
	if err != nil {
		return err
	}
	return nil
}
*/

func (this *Device) Clear(cascadeId int, redraw bool) error {
	if cascadeId >= 0 {
		for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
			this.buffer[cascadeId*MAX7219_DIGIT_COUNT+i] = 0x00
		}
	} else {
		for i := 0; i < this.cascaded*MAX7219_DIGIT_COUNT; i++ {
			this.buffer[i] = 0x00
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

/*
func (this *Device) ClearAll(redraw bool) error {
	for i := 0; i < this.cascaded; i++ {
		err := this.Clear(i, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

*/

func (this *Device) ClearAll(redraw bool) error {
	for i := 0; i < this.cascaded; i++ {
		err := this.Clear(i, false)
		if err != nil {
			return err
		}
	}
	if redraw {
		this.Command(MAX7219_REG_SHUTDOWN, OP_OFF) // not shutdown mode
		this.Command(MAX7219_REG_DISPLAYTEST, 0)   // no display test
		this.Command(MAX7219_REG_SCANLIMIT, 7)     // show all 8 digits
		this.Command(MAX7219_REG_DECODEMODE, 0)    // use individual segments
		this.Command(MAX7219_REG_SHUTDOWN, OP_ON)  // not shutdown mode
		this.Command(MAX7219_REG_NOOP, 0)
		err := this.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) InitRegisters() {
	//this.Command(MAX7219_REG_SHUTDOWN, OP_OFF) // not shutdown mode
	//this.Command(MAX7219_REG_DISPLAYTEST, 0)  // no display test
	//this.Command(MAX7219_REG_SCANLIMIT, 7)    // show all 8 digits
	//this.Command(MAX7219_REG_DECODEMODE, 0) // use individual segments
	//this.Command(MAX7219_REG_SHUTDOWN, OP_ON) // not shutdown mode
	//this.Command(MAX7219_REG_NOOP, 0)
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
			this.Command(Max7219Reg(len(toDisplay)-i), 119) //A
		case "B", "b":
			this.Command(Max7219Reg(len(toDisplay)-i), 31) //b
		case "C":
			this.Command(Max7219Reg(len(toDisplay)-i), 78) //C
		case "c":
			this.Command(Max7219Reg(len(toDisplay)-i), 13) //c
		case "D", "d":
			this.Command(Max7219Reg(len(toDisplay)-i), 61) //d
		case "E", "e":
			this.Command(Max7219Reg(len(toDisplay)-i), 79) //E
		case "F", "f":
			this.Command(Max7219Reg(len(toDisplay)-i), 71) //F
		case "G", "g":
			this.Command(Max7219Reg(len(toDisplay)-i), 94) //G
		case "H":
			this.Command(Max7219Reg(len(toDisplay)-i), 55) //H
		case "h":
			this.Command(Max7219Reg(len(toDisplay)-i), 23) //h
		case "I", "i":
			this.Command(Max7219Reg(len(toDisplay)-i), 6) //i
		case "L", "l":
			this.Command(Max7219Reg(len(toDisplay)-i), 14) //L
		case "N", "n":
			this.Command(Max7219Reg(len(toDisplay)-i), 21) //n
		case "O":
			this.Command(Max7219Reg(len(toDisplay)-i), 99) //O
		case "o":
			this.Command(Max7219Reg(len(toDisplay)-i), 29) //o
		case "P", "p":
			this.Command(Max7219Reg(len(toDisplay)-i), 103) //p
		case "R", "r":
			this.Command(Max7219Reg(len(toDisplay)-i), 5) //r
		case "S", "s":
			this.Command(Max7219Reg(len(toDisplay)-i), 91) //S
		case "T", "t":
			this.Command(Max7219Reg(len(toDisplay)-i), 15) //t
		case "U":
			this.Command(Max7219Reg(len(toDisplay)-i), 62) //U
		case "u":
			this.Command(Max7219Reg(len(toDisplay)-i), 28) //u
		case "Y", "y":
			this.Command(Max7219Reg(len(toDisplay)-i), 59) //Y
		case "0":
			this.Command(Max7219Reg(len(toDisplay)-i), 126) //0
		case "1":
			this.Command(Max7219Reg(len(toDisplay)-i), 48) //1
		case "2":
			this.Command(Max7219Reg(len(toDisplay)-i), 109) //2
		case "3":
			this.Command(Max7219Reg(len(toDisplay)-i), 121) //3
		case "4":
			this.Command(Max7219Reg(len(toDisplay)-i), 51) //4
		case "5":
			this.Command(Max7219Reg(len(toDisplay)-i), 91) //5
		case "6":
			this.Command(Max7219Reg(len(toDisplay)-i), 95) //6
		case "7":
			this.Command(Max7219Reg(len(toDisplay)-i), 112) //7
		case "8":
			this.Command(Max7219Reg(len(toDisplay)-i), 127) //8
		case "9":
			this.Command(Max7219Reg(len(toDisplay)-i), 115) //9
		case ".":
			this.Command(Max7219Reg(len(toDisplay)-i), 128) //.
		case "-":
			this.Command(Max7219Reg(len(toDisplay)-i), 129) //-
		case "=":
			this.Command(Max7219Reg(len(toDisplay)-i), 9) //= lower
		case "J", "j":
			this.Command(Max7219Reg(len(toDisplay)-i), 64) //= upper
		case "K", "k":
			this.Command(Max7219Reg(len(toDisplay)-i), 72) //= split
		case "M", "m":
			this.Command(Max7219Reg(len(toDisplay)-i), 73) //==

		}
	}
	return nil
}
