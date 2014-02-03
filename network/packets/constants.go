package packets

import "github.com/tobz/phosphorus/network"

const (
	RequestBadNameCheck  PacketCode = 0x6A ^ 168
    ResponseBadNameCheck PacketCode = 0xC3
)
