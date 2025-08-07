# renogy-modbus-mqtt
Utility for reading Renogy Rover charge controller data and publishing it to
other systems. Data can be sent to a generic MQTT broker or directly into the
Victron Venus OS dbus where it will appear as a `com.victronenergy.solarcharger`
service.

## Venus OS

To publish values on a Venus device set the environment variable
`RMM_PUSH_MODE=dbus`. The application will register on the system dbus and
expose key solar metrics such as battery voltage, charge current and PV input
voltage.
