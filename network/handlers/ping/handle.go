package ping

import (
	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestPing, HandlePingRequest)
}

func HandlePingRequest(c interfaces.Client, p *network.InboundPacket) error {
	// Not sure what these are for.
	p.Skip(4)

	// Set the ping time to now.
	c.MarkPingTime()

	// Pull out the timestamp so we can send our response.
	timestamp, err := p.ReadUInt32()
	if err != nil {
		return err
	}

	return SendPingResponse(c, timestamp, p.Sequence)
}

func SendPingResponse(c interfaces.Client, timestamp uint32, sequence uint16) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.ResponsePing)

	p.WriteUInt32(timestamp)
	p.WriteRepeated(0x00, 4)
	p.WriteUInt16(sequence + 1)
	p.WriteRepeated(0x00, 6)

	return c.Send(p)
}
