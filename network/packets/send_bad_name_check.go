package packets

import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/network"

func SendBadNameCheckResponse(client interfaces.Client, characterName string, validName bool) error {
    // Set up our outbound packet.
    packet := network.NewOutboundPacket(network.PacketType_TCP, 0xC3)

    // Stick the character name in there so the client knows what we're responding about.
    packet.WriteBoundedString(characterName, 30)

    // They should have an account at this point, but this is just for correctness.
    if client.Account() != nil {
        packet.WriteBoundedString(client.Account().Name(), 20)
    } else {
        packet.WriteRepeated(0x00, 20)
    }

    // Tell them if the name was valid or not.
    if validName {
        packet.WriteUInt8(0x01)
    } else {
        packet.WriteUInt8(0x00)
    }

    return client.Send(packet)
}
