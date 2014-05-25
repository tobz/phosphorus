package network

import "io"
import "fmt"

import "github.com/rcrowley/go-metrics"

import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/utils"

type PacketReader struct {
	conn    io.Reader
	readBuf []byte
	readOff int

	processedBytes metrics.Counter
}

func NewPacketReader(conn io.Reader, counter metrics.Counter) *PacketReader {
	return &PacketReader{conn, make([]byte, 8192), 0, counter}
}

func (r *PacketReader) Next() (*InboundPacket, error) {
	// See if we have any packets in our left over buffer.
	bufferedPacket, err := r.getPacketFromBuffer()
	if bufferedPacket != nil || err != nil {
		return bufferedPacket, err
	}

	// We've cleared through any buffered packets, so do a network read to get everything off the wire.
	n, err := r.conn.Read(r.readBuf[r.readOff:])
	if err != nil {
		return nil, err
	}

	// If we didn't get anything, then we won't be able to get a packet out.
	if n == 0 {
		return nil, nil
	}

	if r.processedBytes != nil {
		r.processedBytes.Inc(int64(n))
	}

	// Increment our read offset based on what we received.
	r.readOff += n

	// See if we can get a packet out yet.
	return r.getPacketFromBuffer()
}

func (r *PacketReader) getPacketFromBuffer() (*InboundPacket, error) {
	// Make sure we have a packet in the buffer to read.
	if !r.hasBufferedPacket() {
		return nil, nil
	}

	// Looks like we have a packet: make sure the packet checksums match.
	err := r.ensurePacketChecksum()
	if err != nil {
		return nil, err
	}

	// Create the packet, keeping in mind to not copy in the checksum.  No need for it now.
	packetLength := r.getPacketLengthFromBuffer()
	packet := NewInboundPacket(r.readBuf[:packetLength-2], constants.PacketTCP)

	// If we pulled out a packet and there's left over data, we need to shift it to the front so the next call
	// will read it out immediately.  Otherwise, we're done here, so put the read offset back to zero to start
	// reading into our buffer from the front.
	if r.readOff > packetLength {
		remaining := (r.readOff - packetLength)
		copy(r.readBuf, r.readBuf[packetLength:])
		r.readOff = remaining
	} else {
		r.readOff = 0
	}

	return packet, nil
}

func (r *PacketReader) hasBufferedPacket() bool {
	// See if we have enough bytes in our buffer for a packet at all.
	if r.readOff < constants.InboundPacketHeaderSize {
		return false
	}

	// See if we have enough bytes for the packet in our buffer based on the header.  Packet length is minus the
	// packet header size, which is why we add it back in.
	packetLength := r.getPacketLengthFromBuffer()

	if r.readOff < packetLength {
		return false
	}

	return true
}

func (r *PacketReader) getPacketLengthFromBuffer() int {
	return ((int(r.readBuf[0]) << 8) | int(r.readBuf[1])) + constants.InboundPacketHeaderSize
}

func (r *PacketReader) ensurePacketChecksum() error {
	packetLength := r.getPacketLengthFromBuffer()

	calculatedChecksum := utils.CalculatePacketChecksum(r.readBuf, 0, packetLength-2)
	providedChecksum := uint16(r.readBuf[packetLength-2])<<8 | uint16(r.readBuf[packetLength-1])
	if providedChecksum != calculatedChecksum {
		return fmt.Errorf("Bad packet checksum: got 0x%X, calculated 0x%X", providedChecksum, calculatedChecksum)
	}

	return nil
}

type PacketWriter struct {
	conn io.Writer

	processedBytes metrics.Counter
}

func NewPacketWriter(conn io.Writer, counter metrics.Counter) *PacketWriter {
	return &PacketWriter{conn, counter}
}

func (w *PacketWriter) Write(p interfaces.Packet) error {
	// Finalize the packet.
	p.Finalize()

	if w.processedBytes != nil {
		w.processedBytes.Inc(int64(p.Len()))
	}

	// Write the packet.
	_, err := w.conn.Write(p.Buffer())
	return err
}
