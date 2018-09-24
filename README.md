# Go package for VE.Direct
go-vedirect can read data from the VE.Direct serial protocol used in compatible products by Victron Energy.

It currently uses the [tarm/serial](https://github.com/tarm/serial) Go package for setting up the serial port.

Has been tested on a Raspberry Pi Zero W connected to a Victron Energy SmartSolar MPPT 75/15.

## Disclaimer
This is my first Go program. It's not pretty (yet). Suggestions for improvements will be highly appreciated.

## Usage

```
s := vedirect.NewStream("/dev/ttyUSB0")
fmt.Println(s)
for {
	block, checksum := s.ReadBlock()
	if checksum == 0 {
		fmt.Println(block)
	} else {
		fmt.Println("Bad block, skipping:", block)
	}
}
```

## TODO

- [ ] Use a larger buffer instead of doing a read() for every byte.
