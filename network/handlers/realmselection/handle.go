package realmselection

import (
	"fmt"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.ClientOnlyRealm, HandleRealmSelection)
}

func HandleRealmSelection(c interfaces.Client, p *network.InboundPacket) error {
	// Not sure what the first byte is.
	p.Skip(1)

	// Get the realm selection.  If this is 0, I believe it indicates the client skipped the realm selection
	// page which is the case when we tell the client that they belong to a specific realm.  Otherwise, it
	// should be a value that matches up with our realm constants.
	clientRealm, err := p.ReadUInt8()
	if err != nil {
		return fmt.Errorf("unable to read realm value")
	}

	// Two more useless bytes.
	p.Skip(2)

	c.Logger().Debug("realmselection", "Client reported realm as %d", clientRealm)

	return nil
}
