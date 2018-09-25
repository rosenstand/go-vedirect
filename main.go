package main

import (
  "flag"
  "fmt"
//  "vedirect"
  "github.com/rosenstand/go-vedirect/vedirect"
)

func main() {
  deviceFlag := flag.String("dev", "/dev/ttyUSB0", "full path to serial device node")
  flag.Parse()
	s := vedirect.NewStream(*deviceFlag)
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
