package givenergy

import (
	"bytes"
	"testing"
)

// TestCRC16Modbus pins the CRC implementation against a known-good value.
// dewet22/givenergy-modbus documents a confirmed real captured request,
// ReadHolding(0x11,0,60), with CRC check=0x474b (see its
// pdu/transparent.py _update_check_code docstring). That request's CRC is
// computed over [addr, func, base_hi, base_lo, count_hi, count_lo] =
// [0x11, 0x03, 0x00, 0x00, 0x00, 0x3C], with the wire bytes being the raw
// CRC's low byte then high byte. Reusing that exact byte sequence here
// (function 0x03, not our 0x04) lets us check our implementation against a
// value from the real device rather than only against itself.
func TestCRC16Modbus(t *testing.T) {
	input := []byte{0x11, 0x03, 0x00, 0x00, 0x00, 0x3C}
	raw := crc16Modbus(input)
	wire := []byte{byte(raw), byte(raw >> 8)}
	want := []byte{0x47, 0x4b}
	if !bytes.Equal(wire, want) {
		t.Fatalf("crc16Modbus(%x) wire bytes = %x, want %x", input, wire, want)
	}
}

func TestBuildReadInputRegistersRequest(t *testing.T) {
	frame := buildReadInputRegistersRequest(0x11, 0, 60)

	// MBAP-style header.
	if got, want := frame[0:2], []byte{0x59, 0x59}; !bytes.Equal(got, want) {
		t.Errorf("tid = %x, want %x", got, want)
	}
	if got, want := frame[2:4], []byte{0x00, 0x01}; !bytes.Equal(got, want) {
		t.Errorf("pid = %x, want %x", got, want)
	}
	if got, want := frame[4:6], []byte{0x00, 0x1c}; !bytes.Equal(got, want) { // len(inner)=26 + 2 = 28 = 0x1c
		t.Errorf("len = %x, want %x", got, want)
	}
	if got, want := frame[6], byte(0x01); got != want {
		t.Errorf("uid = %x, want %x", got, want)
	}
	if got, want := frame[7], byte(0x02); got != want {
		t.Errorf("main function = %x, want %x", got, want)
	}

	if got, want := len(frame), 34; got != want {
		t.Fatalf("frame length = %d, want %d", got, want)
	}

	// device address + transparent function code, at offset 8+10+8=26.
	if got, want := frame[26], byte(0x11); got != want {
		t.Errorf("device address = %x, want %x", got, want)
	}
	if got, want := frame[27], byte(0x04); got != want {
		t.Errorf("transparent function = %x, want %x", got, want)
	}
	// base register (2B) + count (2B).
	if got, want := frame[28:32], []byte{0x00, 0x00, 0x00, 0x3c}; !bytes.Equal(got, want) {
		t.Errorf("base/count = %x, want %x", got, want)
	}
}

func TestParseTransparentResponseRoundTrip(t *testing.T) {
	// Build a synthetic response frame matching the real layout: mbap(8) +
	// adapter_serial(10) + padding(8) + addr(1) + func(1) + inverter_serial(10)
	// + base(2) + count(2) + values(count*2) + crc(2).
	values := make([]byte, 4)
	values[0], values[1] = 0x00, 0x2a // register 0 = 42
	values[2], values[3] = 0x00, 0x64 // register 1 = 100

	inner := []byte{}
	inner = append(inner, padString("AB1234G567", 10)...)
	inner = append(inner, 0, 0, 0, 0, 0, 0, 0, 0x8A) // padding = 0x8A
	inner = append(inner, 0x11, 0x04)                // addr, transparent func (no error bit)
	inner = append(inner, padString("SA1234G567", 10)...)
	inner = append(inner, 0x00, 0x3B) // base register = 59
	inner = append(inner, 0x00, 0x02) // count = 2
	inner = append(inner, values...)

	crcInput := inner[18:] // device-address byte onward
	raw := crc16Modbus(crcInput)
	inner = append(inner, byte(raw), byte(raw>>8))

	frame := []byte{}
	frame = append(frame, 0x59, 0x59, 0x00, 0x01)
	frame = append(frame, byte(len(inner)+2>>8), byte(len(inner)+2))
	frame = append(frame, 0x01, 0x02)
	frame = append(frame, inner...)

	resp, err := parseTransparentResponse(frame)
	if err != nil {
		t.Fatalf("parseTransparentResponse: %v", err)
	}
	if !resp.crcOK {
		t.Errorf("crcOK = false, want true")
	}
	if resp.baseRegister != 59 {
		t.Errorf("baseRegister = %d, want 59", resp.baseRegister)
	}
	if len(resp.values) != 2 || resp.values[0] != 42 || resp.values[1] != 100 {
		t.Errorf("values = %v, want [42 100]", resp.values)
	}
}
