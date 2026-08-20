// Package givenergy implements just enough of GivEnergy's proprietary
// Modbus-TCP dialect to read live telemetry from an inverter's WiFi/ethernet
// dongle. The wire protocol (frame layout, CRC scheme, register map) was
// reverse-engineered by the dewet22/givenergy-modbus project
// (https://github.com/dewet22/givenergy-modbus, Apache License 2.0) and is
// reimplemented here from scratch in Go; see NOTICE.md for attribution.
package givenergy

import "encoding/binary"

const (
	headerTID           = 0x5959 // fixed transaction id GivEnergy always uses
	headerPID           = 0x0001
	mainFuncHeartbeat   = 0x01
	mainFuncTransparent = 0x02

	transparentFuncReadInput = 0x04
	transparentErrorFlag     = 0x80

	requestPadding = 0x08 // GivEnergy's "transparent" sub-frame padding field for requests

	// dataAdapterSerial is sent in the request's data-adapter-serial field.
	// It's ignored by the dongle for client-originated requests.
	dataAdapterSerial = "AB1234G567"
)

// crc16Modbus implements the standard CRC-16/MODBUS algorithm (poly 0xA001,
// init 0xFFFF), which GivEnergy uses to check its request/response frames.
func crc16Modbus(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b)
		for range 8 {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// appendCRC computes the CRC-16/MODBUS over crcInput and appends its two
// wire bytes (low byte first) to frame, matching GivEnergy's convention.
func appendCRC(frame, crcInput []byte) []byte {
	raw := crc16Modbus(crcInput)
	return append(frame, byte(raw), byte(raw>>8))
}

func padString(s string, n int) []byte {
	b := []byte(s)
	if len(b) >= n {
		return b[len(b)-n:]
	}
	out := make([]byte, 0, n)
	for i := 0; i < n-len(b); i++ {
		out = append(out, '*')
	}
	return append(out, b...)
}

// buildReadInputRegistersRequest encodes a full GivEnergy "Transparent"
// request frame asking device addr for count input registers starting at
// base. addr is typically 0x11 for the inverter itself.
func buildReadInputRegistersRequest(addr byte, base, count uint16) []byte {
	inner := make([]byte, 0, 26)
	inner = append(inner, padString(dataAdapterSerial, 10)...)
	inner = binary.BigEndian.AppendUint64(inner, requestPadding)
	inner = append(inner, addr, transparentFuncReadInput)
	inner = binary.BigEndian.AppendUint16(inner, base)
	inner = binary.BigEndian.AppendUint16(inner, count)

	crcInput := []byte{addr, transparentFuncReadInput}
	crcInput = binary.BigEndian.AppendUint16(crcInput, base)
	crcInput = binary.BigEndian.AppendUint16(crcInput, count)
	inner = appendCRC(inner, crcInput)

	frame := make([]byte, 0, 8+len(inner))
	frame = binary.BigEndian.AppendUint16(frame, headerTID)
	frame = binary.BigEndian.AppendUint16(frame, headerPID)
	frame = binary.BigEndian.AppendUint16(frame, uint16(len(inner)+2))
	frame = append(frame, 0x01, mainFuncTransparent)
	frame = append(frame, inner...)
	return frame
}
