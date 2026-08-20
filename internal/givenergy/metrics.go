package givenergy

// Snapshot is the subset of GivEnergy telemetry this tool cares about, all
// read from a single "read input registers 0-59" request.
type Snapshot struct {
	BatterySOC   uint16 // percent, 0-100
	SolarWatts   int    // PV string 1 + string 2 combined
	GridWatts    int    // signed: positive = exporting to grid, negative = importing
	LoadWatts    int    // house consumption
	BatteryWatts int    // signed: positive = charging, negative = discharging
}

// Register offsets within the IR(0..59) block, confirmed against
// givenergy_modbus/model/inverter.py (dewet22/givenergy-modbus, Apache-2.0).
const (
	regPPV1        = 18
	regPPV2        = 20
	regPGridOut    = 30
	regPLoadDemand = 42
	regPBattery    = 52
	regBatterySOC  = 59
)

func snapshotFromRegisters(v []uint16) Snapshot {
	get := func(i int) uint16 {
		if i < 0 || i >= len(v) {
			return 0
		}
		return v[i]
	}
	signed := func(i int) int {
		return int(int16(get(i)))
	}
	return Snapshot{
		BatterySOC:   get(regBatterySOC),
		SolarWatts:   int(get(regPPV1)) + int(get(regPPV2)),
		GridWatts:    signed(regPGridOut),
		LoadWatts:    int(get(regPLoadDemand)),
		BatteryWatts: -signed(regPBattery), // register is negative while charging; negate so field matches its doc comment
	}
}
