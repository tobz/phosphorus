package packets

import "strings"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/network"
import "github.com/tobz/phosphorus/managers"

func init() {
    managers.DefaultPacketManager.RegisterRequestHandler(network.PacketType_TCP, 0x6A ^ 168, HandleBadNameCheckRequest)
}

func HandleBadNameCheckRequest(client *interfaces.Client, packet *network.InboundPacket) error {
    characterName := packet.ReadString(30)

    validName := true

    // See if any of the invalid words we have listed are in this character name.
    for _, invalidWord := managers.DefaultNameManager.InvalidWords {
        if strings.Contains(characterName, invalidWord) {
            validName = false
            break
        }
    }

    return SendBadNameCheckResponse(client, characterName, validName)
}
