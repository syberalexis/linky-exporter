package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type LinkyConnector struct {
	Mode      LinkyMode
	Device    string
	BaudRate  int
	FrameSize int
	Parity    serial.Parity
	StopBits  serial.StopBits
}

// Detect serial connection mode
func (connector *LinkyConnector) Detect() error {
	log.Info("Trying to auto detect TIC mode...")

	if connector.trySerial(Standard) {
		log.Info("Standard Mode detected !")
		connector.Mode = Standard
		connector.BaudRate = Standard.BaudRate
		connector.FrameSize = Standard.FrameSize
		connector.Parity = Standard.Parity
		connector.StopBits = Standard.StopBits
		return nil
	} else {
		log.Debug("It's not standard mode !")
	}

	if connector.trySerial(Historical) {
		log.Info("Historical Mode detected !")
		connector.Mode = Historical
		connector.BaudRate = Historical.BaudRate
		connector.FrameSize = Historical.FrameSize
		connector.Parity = Historical.Parity
		connector.StopBits = Historical.StopBits
		return nil
	} else {
		log.Debug("It's not historical mode !")
	}

	return fmt.Errorf("Impossible to auto detect TIC mode !")
}

// Try serial connection and reading
func (connector LinkyConnector) trySerial(mode LinkyMode) bool {
	m := &serial.Mode{BaudRate: mode.BaudRate, DataBits: mode.FrameSize, Parity: mode.Parity, StopBits: mode.StopBits}
	stream, err := serial.Open(connector.Device, m)
	if err != nil {
		return false
	}

	reader := bufio.NewReader(stream)
	regex, _ := regexp.Compile(`^[A-Z0-9\-+]+ +[a-zA-Z0-9 \.\-]+ +.$`)

	log.Debug("Read serial data...")
	for i := 1; i <= 5; i++ {
		bytes, _, _ := reader.ReadLine()
		line := string(bytes)
		log.Debug("Try line ", i, "/5 : ", line)
		if regex.MatchString(line) {
			return true
		} else {
			log.Debug("Regex not match")
		}
	}
	return false
}

// Read serial values
func (connector LinkyConnector) readSerial() ([][]string, error) {
	log.Debug("Read serial with config device:", connector.Device, " baudrate:", connector.BaudRate, " framesize:", connector.FrameSize, " parity:", connector.Parity, " stopbits:", connector.StopBits)
	m := &serial.Mode{BaudRate: connector.BaudRate, DataBits: connector.FrameSize, Parity: connector.Parity, StopBits: connector.StopBits}
	stream, err := serial.Open(connector.Device, m)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(stream)
	started := false
	var values [][]string

	log.Debug("Read serial data...")
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		line := string(bytes)

		// End loop when block ended
		if started && strings.ContainsRune(line, 0x03) {
			stream.Close()
			break
		}

		// Collect data line by line
		if started {
			log.Debug(line)
			values = append(values, strings.FieldsFunc(line, func(r rune) bool { return r == 0x09 || r == ' ' }))
		}

		// Start reading data when block started
		if strings.ContainsRune(line, 0x02) {
			started = true
		}
	}
	log.Debug("Read serial data ended !")

	return values, nil
}

// Return last serial Historical TIC
func (connector LinkyConnector) GetLastHistoricalTicValue() (*HistoricalTicValue, error) {
	lines, err := connector.readSerial()

	if err != nil {
		log.Errorf("Failed to read historical serial : %s", err)
		return nil, err
	}

	values := HistoricalTicValue{}
	for _, line := range lines {
		values.ParseParam(line[0], line[1:])
	}

	return &values, nil
}

// Return last serial Standard TIC
func (connector LinkyConnector) GetLastStandardTicValue() (*StandardTicValue, error) {
	lines, err := connector.readSerial()

	if err != nil {
		log.Errorf("Failed to read standard serial : %s", err)
		return nil, err
	}

	values := StandardTicValue{}
	for _, line := range lines {
		values.ParseParam(line[0], line[1:])
	}

	return &values, nil
}

// Parse parity from string to serial object
func ParseParity(value string) (parity serial.Parity) {
	switch value {
	case "ParityNone", "N":
		parity = serial.NoParity
		break
	case "ParityOdd", "O":
		parity = serial.OddParity
		break
	case "ParityEven", "E":
		parity = serial.EvenParity
		break
	case "ParityMark", "M":
		parity = serial.MarkParity
		break
	case "ParitySpace", "S":
		parity = serial.SpaceParity
		break
	default:
		log.Error(fmt.Errorf("Impossible to parse Parity named : %s", value))
		os.Exit(3)
	}
	return
}

// Parse stop bits from string to serial object
func ParseStopBits(value string) (stopBits serial.StopBits) {
	switch value {
	case "Stop1", "1":
		stopBits = serial.OneStopBit
		break
	case "Stop1Half", "15":
		stopBits = serial.OnePointFiveStopBits
		break
	case "Stop2", "2":
		stopBits = serial.TwoStopBits
		break
	default:
		log.Error(fmt.Errorf("Impossible to parse StopBits named : %s", value))
		os.Exit(3)
	}
	return
}
