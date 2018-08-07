package mcp9808

import (
	"log"

	i2c "github.com/fedeonline/i2c-go"
)

const (
	// Default I2C address for device.
	defaultI2CAddr = 0x18
	hwBus          = 1 // Depends on Raspberry Pi version (0 or 1)

	// Register addresses.
	regConfig         = 0x01
	regUpperTemp      = 0x02
	regLowerTemp      = 0x03
	regCriticalTemp   = 0x04
	regAmbientTemp    = 0x05
	regManufacturerID = 0x06
	regDeviceID       = 0x07
	regResolution     = 0x08

	// Default values
	chkManufacturerID = 0x54
	chkDeviceID       = 0x0400

	// Configuration register values.
	regCongigShutdown   = 0x0100
	regConfigCritLocked = 0x0080
	regConfigWinLocked  = 0x0040
	regConfigIntClr     = 0x0020
	regConfigAlertStat  = 0x0010
	regConfigAlertCtrl  = 0x0008
	regConfigAlertSel   = 0x0002
	regConfigAlertPol   = 0x0002
	regConfigAlertMode  = 0x0001
)

// Find any modules and return the addresses of the detected
func Find() (m []uint8) {
	for i := uint8(0); i < 8; i++ {
		addr := defaultI2CAddr + i
		discovered := Check(addr)
		if discovered == true {
			log.Printf("MCP9808 detected at address 0x%x", addr)
			m = append(m, addr)
		}
	}
	return m
}

// Check for a valid MCP9808 sensor at the i2c address
func Check(address uint8) bool {
	i, err := i2c.NewI2C(address, hwBus)
	if err != nil {
		return false
	}
	defer i.Close()

	if v, err := i.ReadRegU16BE(regManufacturerID); err != nil || v != chkManufacturerID {
		return false
	}
	if v, err := i.ReadRegU16BE(regDeviceID); err != nil || v != chkDeviceID {
		return false
	}
	return true
}

// ReadAmbientTemp the current temperature from the MCP9808 sensor at the i2c address
func ReadAmbientTemp(address uint8) (float32, error) {
	sensor, err := i2c.NewI2C(address, hwBus)
	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer sensor.Close()

	t, err := sensor.ReadRegU16BE(regAmbientTemp)
	if err != nil {
		return 0, err
	}

	return float32(t&0x0FFF) / float32(16), nil
}
