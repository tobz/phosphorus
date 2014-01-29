package packets

import "github.com/tobz/phosphorus"
import "github.com/tobz/phosphorus/network"

type BadNameCheckResponse struct {
    Name string
    Valid bool
}

func (r *BadNameCheckResponse) GetPacketType() network.PacketType {
    return network.PacketType_TCP
}

func (r *BadNameCheckResponse) GetPacketCode() network.PacketCode {
    return 0xC3
}

func (r *BadNameCheckResponse) SendResponse(client *phosphorus.Client) error {
    packet := network.NewOutboundPacket(r.GetPacketType(), r.GetPacketCode())

    packet.WriteBoundedString(r.Name, 30)

    if client.Account != nil {
        packet.WriteBoundedString(client.Account.Name, 20)
    } else {
        packet.WriteRepeated(0x00, 20)
    }

    packet.WriteUInt8(valid ? 0x01 : 0x00)

    client.Send(packet)
}
