package phosphorus

type BadNameCheckResponse struct {
    Name string
    Valid bool
}

func (r *BadNameCheckResponse) GetPacketType() PacketType {
    return network.PacketType_TCP
}

func (r *BadNameCheckResponse) GetPacketCode() PacketCode {
    return 0xC3
}

func (r *BadNameCheckResponse) SendResponse(client *Client) error {
    packet := NewOutboundPacket(r.GetPacketType(), r.GetPacketCode())

    packet.WriteBoundedString(r.Name, 30)

    if client.Account != nil {
        packet.WriteBoundedString(client.Account.Name, 20)
    } else {
        packet.WriteRepeated(0x00, 20)
    }

    if valid {
        packet.WriteUInt8(0x01)
    } else {
        packet.WriteUInt8(0x00)
    }

    client.Send(packet)
}
