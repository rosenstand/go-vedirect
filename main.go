package main

import (
  "fmt"
//  "vedirect"
  "github.com/rosenstand/go-vedirect/vedirect"
)

func main() {
	s := vedirect.NewStream("/dev/ttyUSB0")
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
