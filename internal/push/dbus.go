package push

import (
	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/dbus"
)

// DBUSPusher publishes controller information onto the Venus OS dbus.
type DBUSPusher struct {
	client *dbus.Client
	inChan chan gorenogymodbus.DynamicControllerInformation
}

// NewDBUSPusher creates a new pusher that writes to dbus.
func NewDBUSPusher(client *dbus.Client, inChan chan gorenogymodbus.DynamicControllerInformation) *DBUSPusher {
	return &DBUSPusher{client: client, inChan: inChan}
}

// Push writes values from the channel to dbus.
func (p *DBUSPusher) Push() error {
	for dci := range p.inChan {
		p.client.Update("/Dc/0/Voltage", dci.BatteryVoltage.InexactFloat64())
		p.client.Update("/Dc/0/Current", dci.ChargingCurrent.InexactFloat64())
		p.client.Update("/Pv/0/Voltage", dci.SolarPanelVoltage.InexactFloat64())
		p.client.Update("/Pv/0/Current", dci.SolarPanelCurrent.InexactFloat64())
		p.client.Update("/Yield/Power", dci.ChargingPower.InexactFloat64())
	}
	return nil
}

// Close closes the underlying dbus connection.
func (p *DBUSPusher) Close() error {
	return p.client.Close()
}
