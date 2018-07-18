// +build linux
// Input device event monitor.
package main

import (
//	"errors"
	"fmt"
	"github.com/gvalkov/golang-evdev"
	"os"
	"strings"
)

const (
	usage       = "usage: evtest <device> [<type> <value>]"
	device_glob = "/dev/input/event*"
)

var choice int
func select_device() (*evdev.InputDevice, error) {
	devices, _ := evdev.ListInputDevices(device_glob)
	lines := make([]string, 0)
	max := 0
	if len(devices) > 0 {
		for i := range devices {
			dev := devices[i]
			str := fmt.Sprintf("%-3d %-20s %-35s %s", i, dev.Fn, dev.Name, dev.Phys)
			if len(str) > max {
				max = len(str)
			}
			lines = append(lines, str)
		}
		fmt.Printf("%-3s %-20s %-35s %s\n", "ID", "Device", "Name", "Phys")
		fmt.Printf(strings.Repeat("-", max) + "\n")
		fmt.Printf(strings.Join(lines, "\n") + "\n")
		choice = 0
		}
	return devices[choice], nil
}

func format_event(ev *evdev.InputEvent) string {
	var res, f, code_name string
	code := int(ev.Code)
	etype := int(ev.Type)
	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			f = " -"
		} else {
			f = " "
		}
		return fmt.Sprintf(f)
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[code]
		if haskey {
			code_name = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				code_name = val
			} else {
				code_name = "?"
			}
		}
	default:
		m, haskey := evdev.ByEventType[etype]
		if haskey {
			code_name = m[code]
		} else {
			code_name = "?"
		}
	}
	evfmt := "shid - "
	res = fmt.Sprintf("%s %s", evfmt, code_name)

	return res
}

func main() {
	var dev *evdev.InputDevice
	var events []evdev.InputEvent
	var err error

	switch len(os.Args) {
	case 1:
		dev, err = select_device()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case 2:
		dev, err = evdev.Open(os.Args[1])
		if err != nil {
			fmt.Printf("unable to open input device: %s\n", os.Args[1])
			os.Exit(1)
		}
	default:
		fmt.Printf(usage + "\n")
		os.Exit(1)
	}

	info := fmt.Sprintf("bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x",
		dev.Bustype, dev.Vendor, dev.Product, dev.Version)

	repeat_info := dev.GetRepeatRate()

	fmt.Printf("Evdev protocol version: %d\n", dev.EvdevVersion)
	fmt.Printf("Device name: %s\n", dev.Name)
	fmt.Printf("Device info: %s\n", info)
	fmt.Printf("Repeat settings: repeat %d. delay %d\n", repeat_info[0], repeat_info[1])
	fmt.Printf("Device capabilities:\n")
	fmt.Printf("shid - listening for events\n")

	for {
		events, err = dev.Read()
		for i := range events {
			str := format_event(&events[i])
			fmt.Println(str)
		}
	}
}
