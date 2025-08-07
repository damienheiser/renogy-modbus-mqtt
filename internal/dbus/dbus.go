package dbus

import (
	"fmt"
	"sync"

	godbus "github.com/godbus/dbus/v5"
)

// Client represents a connection to the system dbus on Venus OS.
type Client struct {
	conn  *godbus.Conn
	items map[string]*busItem
	mu    sync.Mutex
}

// busItem represents a single com.victronenergy.BusItem on the dbus.
type busItem struct {
	conn  *godbus.Conn
	path  godbus.ObjectPath
	value godbus.Variant
}

// GetValue returns the current value stored at this path.
func (b *busItem) GetValue() (godbus.Variant, *godbus.Error) {
	return b.value, nil
}

// SetValue updates the value and emits the ValueChanged signal on dbus.
func (b *busItem) SetValue(v godbus.Variant) *godbus.Error {
	b.value = v
	_ = b.conn.Emit(b.path, "com.victronenergy.BusItem.ValueChanged", v)
	return nil
}

// GetText returns the textual representation of the value.
func (b *busItem) GetText() (string, *godbus.Error) {
	return fmt.Sprint(b.value.Value()), nil
}

// newBusItem creates a dbus object for the given path.
func newBusItem(conn *godbus.Conn, path string) *busItem {
	item := &busItem{conn: conn, path: godbus.ObjectPath(path)}
	conn.Export(item, item.path, "com.victronenergy.BusItem")
	return item
}

// NewClient connects to the system dbus and registers the given service name.
func NewClient(service string) (*Client, error) {
	conn, err := godbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("connect system bus: %w", err)
	}

	reply, err := conn.RequestName(service, godbus.NameFlagDoNotQueue)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("request name: %w", err)
	}
	if reply != godbus.RequestNameReplyPrimaryOwner {
		conn.Close()
		return nil, fmt.Errorf("name %s already taken", service)
	}

	c := &Client{
		conn:  conn,
		items: make(map[string]*busItem),
	}

	// Set some mandatory management items so Venus OS recognises the service.
	c.Update("/Mgmt/ProcessName", "renogy-modbus-mqtt")
	c.Update("/ProductName", "Renogy Rover")
	c.Update("/DeviceInstance", 0)

	return c, nil
}

// Update sets the value for a path on the dbus.
func (c *Client) Update(path string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[path]
	if !ok {
		item = newBusItem(c.conn, path)
		c.items[path] = item
	}

	item.SetValue(godbus.MakeVariant(value))
}

// Close closes the underlying dbus connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
