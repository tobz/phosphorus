package network

import "io"
import "fmt"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/utils"

type PacketReader struct {
	conn    io.Reader
	readBuf []byte
	readOff int
}

func NewPacketReader(conn io.Reader) *PacketReader {
	return &PacketReader{conn, make([]byte, 8192), 0}
}

func (r *PacketReader) Next() (*InboundPacket, error) {
	// do a read to see if we have anything
	n, err := r.conn.Read(r.readBuf[r.readOff:])
	if err != nil {
		return nil, err
	}

	// if we didn't get anything, then we won't be able to get a packet out
	if n == 0 {
		return nil, nil
	}

	// we need atleast InboundPacketHeaderSize bytes to even have a valid empty packet
	bufSize := r.readOff + n
	if bufSize < constants.InboundPacketHeaderSize {
		// increment our offset and dip out
		r.readOff += n
		return nil, nil
	}

	// figure out how long are packet is.  packet length is always length of payload, so we have to
	// add in the header size here to properly grab it all.
	packetLength := ((int(r.readBuf[0]) << 8) | int(r.readBuf[1])) + constants.InboundPacketHeaderSize

	// make sure the packet's checksum is on point.
	calculatedChecksum := utils.CalculatePacketChecksum(r.readBuf, 0, packetLength-2)
	providedChecksum := uint16(r.readBuf[packetLength-2])<<8 | uint16(r.readBuf[packetLength-1])
	if providedChecksum != calculatedChecksum {
		return nil, fmt.Errorf("bad packet checksum: got %X, calculated %X", providedChecksum, calculatedChecksum)
	}

	// create the packet, skipping the checksum at the end
	packet := NewInboundPacket(r.readBuf[:packetLength-2], constants.PacketTCP)

	// shift any remaining data to the beginning of the read buf.
	if bufSize > packetLength {
		copy(r.readBuf, r.readBuf[packetLength:])
	}

	// reset our read buffer to try and red in another packet next time around
	r.readOff = 0

	return packet, nil
}

type PacketWriter struct {
	conn io.Writer
}

func NewPacketWriter(conn io.Writer) *PacketWriter {
	return &PacketWriter{conn}
}

func (r *PacketWriter) Write(p interfaces.Packet) error {
	// Finalize the packet.
    p.Finalize()

    // Write the packet.
    _, err := r.conn.Write(p.Buffer())
    return err
}
