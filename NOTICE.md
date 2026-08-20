# Notices

Copyright (c) 2026 Dan Dukeson <dandukeson@gmail.com>. Licensed under the
MIT License (see [LICENSE](LICENSE)).

## Third-party attribution

This project's understanding of GivEnergy's proprietary Modbus-TCP dialect
(frame layout, CRC scheme, and register map) comes from reading the source of:

**[dewet22/givenergy-modbus](https://github.com/dewet22/givenergy-modbus)**
Licensed under the Apache License, Version 2.0
(https://www.apache.org/licenses/LICENSE-2.0)

No code from that project is copied here — this is a from-scratch Go
implementation — but the protocol knowledge (frame structure in
`framer.py`/`codec.py`/`pdu/transparent.py`/`pdu/read_registers.py`, the
register offsets in `model/inverter.py`, and one real captured request's CRC
value used as a test vector) is derived from that Apache-2.0-licensed work,
and is credited here per its license terms.
