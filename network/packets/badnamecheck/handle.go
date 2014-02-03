package badnamecheck

import (
    "strings"

    "github.com/tobz/phosphorus/network"
    "github.com/tobz/phosphorus/network/handlers"
    "github.com/tobz/phosphorus/server"
)

func init() {
    handlers.Register(network.PacketTCP, RequestBadNameCheck, HandleBadNameCheck)
}

func HandleBadNameCheck(client *server.Client, p *network.InboundPacket) error {
    name := p.ReadString(30)
    valid := true

    for _, word := range client.Server.Config.InvalidWords {
        if strings.Contains(name, word) {
            valid = false
            break
        }
    }

    return SendBadNameCheckResponse(client, name, valid)
}

func SendBadNameCheckResponse(client *server.Client, name string, isValid bool) error {
    p := network.NewOutboundPacket(network.PacketTCP, ResponseBadNameCheck)
    p.WriteBoundedString(name, 30)

    if client.Account != nil {
        panic("client didn't have an account!")
    }

    p.WriteRepeated(0x00, 20)
    if isValid {
        p.WriteUint8(0x01)
    } else {
        p.WriteUint8(0x00)
    }

    return client.Send(p)
}
