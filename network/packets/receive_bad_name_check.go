package phosphorus

import "strings"

func init() {
    handler := &BadNameCheckRequest{BaseHandler{
        PacketType: PacketType_TCP,
        PacketCode: 0x6A ^ 168,
    }}

    DefaultPacketManager.RegisterRequestHandler(handler)
}

// Handles checking if the supplied character name contains invalid words, such as racial epithets and expletives.
type BadNameCheckRequest struct {
    BaseHandler
}

func (r *BadNameCheckRequest) HandleRequest(client *Client, packet *InboundPacket) error {
    characterName := packet.ReadString(30)

    validName := true

    // See if any of the invalid words we have listed are in this character name.
    for _, invalidWord := DefaultNameManager.InvalidWords {
        if strings.Contains(characterName, invalidWord) {
            validName = false
            break
        }
    }

    return DefaultPacketManager.Send(client, &BadNameCheckResponse{ Name: characterName, Valid: validName })
}
