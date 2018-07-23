package main

import (
	"fmt"
	"github.com/gvalkov/golang-evdev"
	"os"
)

const (
	device_glob = "/dev/input/event*"
)

func s_device() (*evdev.InputDevice, error) {
	devices, _ := evdev.ListInputDevices(device_glob)

	for _, dev := range devices {
		fmt.Printf("%s %s %s \n", dev.Fn, dev.Name, dev.Phys)
	}
	return devices[10], nil
}

func f_event(ev *evdev.InputEvent) string {
	var res, f, code_name string
	code := int(ev.Code)
	etype := int(ev.Type)
	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			f = "flag to trigger input, aux to golang"
		}
		return fmt.Sprintf("%s", f)//, evdev.SYN[code])
	default:
		m, haskey := evdev.ByEventType[etype]
		if haskey {
			code_name = m[code]
		} else {
			code_name = "?"
		}
	}
	res = fmt.Sprintf("%s", code_name)
	return res
}

func main() {
	var dev *evdev.InputDevice
	var events []evdev.InputEvent
	var err error

	switch len(os.Args){
	case 1:
		dev, err =s_device()
		fmt.Println(err)
	default:
		fmt.Println("todo bien,empecemos lecutra SHID")
	}

	fmt.Println("todo ok, accadiendo a la entrada de lectura")
	info := fmt.Sprintf("bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x",dev.Bustype, dev.Vendor, dev.Product, dev.Version)
	fmt.Printf("Evdev protocol version: %d\n", dev.EvdevVersion)
	fmt.Printf("info %s Device name: %s\n",info, dev.Name)
	fmt.Printf("Listening for events ...\n")
	for {
		events, err = dev.Read()
		for i := range events {
			str := f_event(&events[i])
			fmt.Println(str)
		}
	}
}
