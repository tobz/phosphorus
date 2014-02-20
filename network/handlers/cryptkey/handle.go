package cryptkey

import (
	"fmt"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestCryptKey, HandleCryptKeyRequest)
}

func HandleCryptKeyRequest(c interfaces.Client, p *network.InboundPacket) error {
	// Get the client type and the encryption stuff.
	useRC4, err := p.ReadUInt8()
	if err != nil {
		return err
	}

	// Client metadata. High bits are client type and lower bits are client addons.
	_, err = p.ReadUInt8()
	if err != nil {
		return err
	}

	// Pull in the client version - three bytes.
	versionMajor, err := p.ReadUInt8()
	if err != nil {
		return err
	}

	versionMinor, err := p.ReadUInt8()
	if err != nil {
		return err
	}

	versionBuild, err := p.ReadUInt8()
	if err != nil {
		return err
	}

	// Now get the right client version value.  We add 900 to the value if it's over 200 as
	// that corresponds to the jump in client versions starting at 1.110.
	versionNumeric := uint16(versionMajor*100) + uint16(versionMinor*10) + uint16(versionBuild)
	if versionNumeric >= 200 {
		versionNumeric += 900
	}

	// Make sure it's a valid version.
	version := constants.ClientVersion(versionNumeric)
	if version < constants.ClientVersionMinimum || version > constants.ClientVersionMaximum {
		return fmt.Errorf("client version reported as %d: allowed range %d to %d", version, constants.ClientVersionMinimum, constants.ClientVersionMaximum)
	}

	c.SetClientVersion(version)
	c.Logger().Debug("cryptkey", "client version: %d", version)

	// If 'useRC4' is == 1, that means there's an RC4 sbox chunk waiting for us.  We have to read it, store it, and then
	// start encrypting the rest of our communications using RC4 seeded from the given sbox.
	if useRC4 == 1 {
		c.Logger().Warn("cryptkey", "client is trying to initiate RC4 encryption, but we don't support it.. yet?")
		return nil
	}

	// No RC4 means that we get to dictate whether or not we want encryption.
	return SendCryptKeyResponse(c)
}

func SendCryptKeyResponse(c interfaces.Client) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.ResponseCryptKey)

	// Send the client version back.
	versionDivisor := 100
	if c.ClientVersion() > 199 {
		versionDivisor = 1000
	}

	p.WriteUInt8(uint8(int(c.ClientVersion()) / versionDivisor))
	p.WriteUInt8(uint8((int(c.ClientVersion()) % versionDivisor) / 10))

	// I think this used to be/was going to be the build component of the client version?  Or
	// it's the "use / don't use" encryption value now.  Not sure yet.
	p.WriteUInt8(0x00)

	return c.Send(p)
}
