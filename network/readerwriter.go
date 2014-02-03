package network

import "io"

type PacketReader struct {
    conn io.Reader
}

func NewPacketReader(conn io.Reader) *PacketReader {
    return &PacketReader{conn}
}

func (r *PacketReader) Next() (*InboundPacket, error) {
    // ...
}

func (r *PacketReader) Stop() {}

type Writer struct {
    conn io.Writer
}

func NewPacketWriter(conn io.Writer) *PacketWriter {
    return &PacketWriter{conn}
}

func (r *PacketWriter) Write(p Packet) error {
    // ...
}

func (r *PacketWriter) Stop() {}
