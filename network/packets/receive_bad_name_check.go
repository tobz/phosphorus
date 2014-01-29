package packets

import "strings"
import "github.com/tobz/phosphorus"
import "github.com/tobz/phosphorus/network"

func init() {
    handler := &BadNameCheckRequest{
        network.BaseHandler{
            PacketType: PacketType_TCP,
            PacketCode: 0x6A ^ 168,
        }
    }

    network.DefaultPacketManager.RegisterRequestHandler(handler)
}

// Handles checking if the supplied character name contains invalid words, such as racial epithets and expletives.
type BadNameCheckRequest struct {
    BaseHandler
}

func (r *BadNameCheckRequest) HandleRequest(client *phosphorus.Client, packet *phosphorus.Packet) error {
    characterName := packet.ReadString(30)

    validName := true

    // See if any of the invalid words we have listed are in this character name.
    for _, invalidWord := phosphorus.DefaultNameManager.InvalidWords {
        if strings.Contains(characterName, invalidWord) {
            validName = false
            break
        }
    }

    return network.DefaultPacketManager.Send(client, &BadNameCheckResponse{ Name: characterName, Valid: validName })
}
