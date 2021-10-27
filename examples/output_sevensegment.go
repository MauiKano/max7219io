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


package main

import (
	"log"
	"github.com/max7219"
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
