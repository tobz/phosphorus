package interfaces

type Client interface {
    Account() Account

    Send(packet Packet) error
}
