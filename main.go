package main

import (
  "flag"
  "fmt"
//  "vedirect"
  "github.com/rosenstand/go-vedirect/vedirect"
)

func main() {
  var device string
  flag.StringVar(&device, "dev", "/dev/ttyUSB0", "full path to serial device node")
  flag.Parse()
	s := vedirect.NewStream(device)
	fmt.Println(s)
	for {
		b, checksum := s.ReadBlock()
		if checksum == 0 {
			fmt.Println(b)
		} else {
			fmt.Println("Bad block, skipping:", b)
		}
	}
}
