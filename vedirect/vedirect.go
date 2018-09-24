package vedirect

import (
	//"bytes"
	"fmt"
	"log"
	//"os"
	//"strings"
	"github.com/tarm/serial"
)

const (
	InChecksum = 1
	InFrame    = 2
	InLabel    = 3
	InValue    = 4
	WaitHeader = 5
)

type Block struct {
	checksum int
	fields   map[string]string
}

func (b Block) Validate() bool {
	if b.checksum%256 == 0 {
		return true
	}
	return false
}

type Stream struct {
	Device string
	//Port   *os.File
	Port  *serial.Port
	State int
}

type Streamer interface {
	Read() int
}

func NewStream(dev string) Stream {
	s := Stream{}
	s.Device = dev
	s.State = 0
	var err error

	c := &serial.Config{Name: s.Device, Baud: 19200}
	s.Port, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	/*
		s.Port, err = os.Open(s.Device)
		if err != nil {
			log.Fatal(err)
		}
	*/
	fmt.Println("Stream initialized:", s)
	return s
}

// Field format: <Newline><Field-Label><Tab><Field-Value>
// Last field in block will always be "Checksum".
// The value is a single byte, and the modulo 256 sum
// of all bytes in a block will equal 0 if there were
// no transmission errors.

func (s *Stream) ReadBlock() (Block, int) {
	var b = Block{}
	b.fields = make(map[string]string)
	var frame_length int = 0
	var prev_state int
	//var label string
	//var value string
	var label = make([]byte, 0, 9)  // VE recommended buffer size.
	var value = make([]byte, 0, 33) // VE recommended buffer size.

	buf := make([]byte, 1)

	for {
		n, err := s.Port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		str := string(buf[:n])
		var char byte = buf[0]

		// HEX mode is documented in BlueSolar-HEX-protocol-MPPT.pdf.
		// catch and ignore VE.Direct HEX frames from stream, otherwise
		// they mess up our checksum and we lose the current block.
		if char == ':' && s.State != InChecksum { // ":": beginning of frame
			//if str == ":" { // ":": beginning of frame
			prev_state = s.State // save state
			s.State = InFrame
			frame_length = 1
			continue
		}
		if s.State == InFrame {
			frame_length = frame_length + 1
			if str == "\n" { // end of frame
				s.State = prev_state // restore state
				fmt.Printf("%d bytes HEX frame ignored\n", frame_length)
			}
			continue // ignore frame contents
		}

		// convert byte to integer and add to sum.
		b.checksum += int(buf[0])

		// end of block. must process before byte evaluation.
		// checksum byte could have any value.
		if s.State == InChecksum {
			s.State = WaitHeader
			//if b.checksum % 256 == 0 {
			return b, b.checksum % 256 // 0 on valid checksum
			//} else {
			//  fmt.Println("Bad block!", b.fields)
			//}
		}

		switch char {
		case 13: // "\r": beginning of field
			if s.State != WaitHeader { // avoid zero-valued entry on first run
				//  b.fields[label] = value
				b.fields[string(label)] = string(value)
			}
			//label = ""
			//value = ""
			label = label[:0] // clear slice
			value = value[:0] // clear slice
			s.State = InLabel
			//continue

		case 10: // "\n": avoid appending \n to label
			//continue

		case 9: // "\t": label/value seperator
			if string(label) == "Checksum" {
				s.State = InChecksum
			} else {
				s.State = InValue
			}
			//continue

		default:
			if s.State == InLabel {
				//label += str
				label = append(label, buf[0])
			} else if s.State == InValue {
				//value += str
				value = append(value, buf[0])
			}
		}
	}
}
