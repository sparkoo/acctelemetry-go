package acctelemetry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const BROADCASTING_PROTOCOL_VERSION byte = 4

const (
	REGISTER_COMMAND_APPLICATION   byte = 1
	UNREGISTER_COMMAND_APPLICATION byte = 9

	// we don't care about below types for reading the telemetry
	REQUEST_ENTRY_LIST           byte = 10
	REQUEST_TRACK_DATA           byte = 11
	CHANGE_HUD_PAGE              byte = 49
	CHANGE_FOCUS                 byte = 50
	INSTANT_REPLAY_REQUEST       byte = 51
	PLAY_MANUAL_REPLAY_HIGHLIGHT byte = 52
	SAVE_MANUAL_REPLAY_HIGHLIGHT byte = 60
)

const (
	REGISTRATION_RESULT byte = 1
	REALTIME_UPDATE     byte = 2
	REALTIME_CAR_UPDATE byte = 3

	// we don't care about below types for reading the telemetry
	ENTRY_LIST         byte = 4
	ENTRY_LIST_CAR     byte = 6
	TRACK_DATA         byte = 5
	BROADCASTING_EVENT byte = 7
)

type connectionResult struct {
	connectionId      int32
	connectionSuccess bool
	readOnly          bool
	errorMessage      string
}

func (telemetry *AccTelemetry) connect() error {
	connectMessage, err := telemetry.createConnectMessage()
	if err != nil {
		return fmt.Errorf("failed to craete connect message: %w", err)
	}

	// send connection request
	_, sendErr := telemetry.udpConnection.Write(connectMessage)
	if sendErr != nil {
		return fmt.Errorf("failed to send connection message: %w", sendErr)
	}

	// give lot of time for connection
	telemetry.udpConnection.SetReadDeadline(time.Now().Add(15 * time.Second))
	inBuffer := make([]byte, 128)
	_, _, err = telemetry.udpConnection.ReadFromUDP(inBuffer)
	if err != nil {
		telemetry.Close()
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("UDP read timeout, ACC probably not running: %w", netErr)
		} else {
			return fmt.Errorf("UDP read failed: %w", err)
		}
	} else {
		fmt.Println("UDP Connected")
		telemetry.realtimeCarUpdate = &RealtimeCarUpdate{
			BestSessionLap: &LapInfo{Splits: [8]int32{}},
			LastLap:        &LapInfo{Splits: [8]int32{}},
			CurrentLap:     &LapInfo{Splits: [8]int32{}},
		}
	}

	connectionResult, err := readConnectionResult(bytes.NewBuffer(inBuffer))
	if err != nil {
		return fmt.Errorf("failed to read connection response: %w", err)
	}
	fmt.Printf("Connected to ACC, listen for messages: '%+v'", connectionResult)

	go func() {
		payload := make([]byte, 128)
		for telemetry.udpConnection != nil {
			// zero payload slice before each use
			// this is to avoid unnecessary memory allocation for each request
			for i := range payload {
				payload[i] = 0
			}

			telemetry.udpConnection.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			_, _, err := telemetry.udpConnection.ReadFromUDP(payload)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					fmt.Printf("UDP read timeout, ACC may not be running: %s", netErr)
				} else {
					fmt.Printf("UDP read failed: %s", err)
				}
			} else {
				if err := telemetry.readMessage(payload); err != nil {
					fmt.Printf("failed to read the message: %s\n", err)
				}
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	return nil
}

func (telemetry *AccTelemetry) createConnectMessage() ([]byte, error) {
	outBuffer := bytes.NewBuffer([]byte{})
	var writeErr error
	writeErr = outBuffer.WriteByte(REGISTER_COMMAND_APPLICATION)
	writeErr = outBuffer.WriteByte(BROADCASTING_PROTOCOL_VERSION)
	writeErr = writeString(outBuffer, telemetry.config.udpDisplayName)
	writeErr = writeString(outBuffer, telemetry.config.UdpConnectionPassword)
	writeErr = binary.Write(outBuffer, binary.LittleEndian, telemetry.config.udpRealtimeUpdateIntervalMS)
	writeErr = writeString(outBuffer, telemetry.config.udpCommandPassword)

	if writeErr != nil {
		return nil, fmt.Errorf("failed to write connection data to byte buffer: %w", writeErr)
	}

	return outBuffer.Bytes(), nil
}

// response format:
// 1 byte - message type
// 4 bytes - int32 connectionId
// 1 byte - connectionSuccess `byte > 0`
// 1 byte - readonly `byte == 0`
// 2 bytes - error length => N
// N bytes - error message
func readConnectionResult(payload *bytes.Buffer) (*connectionResult, error) {
	result := &connectionResult{}

	// read type
	messageType, err := payload.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read response type: %w", err)
	} else if messageType != REGISTRATION_RESULT {
		return nil, fmt.Errorf("should've received REGISTRATION_RESULT['1'] message type, but got '%d'", messageType)
	}

	// read connectionId
	connectionId, err := readInt32(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read connectionId: %w", err)
	}
	result.connectionId = connectionId

	// read success
	connectionSuccess, err := payload.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read conn success: %w", err)
	}
	result.connectionSuccess = connectionSuccess > 0

	// read readonly
	readonly, err := payload.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read readonly: %w", err)
	}
	result.readOnly = readonly == 0

	// read error
	errorLength, err := readUint16(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read error message length: %w", err)
	}

	//TODO: handle reading error message if this is non zero
	if errorLength > 0 {
		result.errorMessage = "some error"
	}

	return result, nil
}

func (t *AccTelemetry) readMessage(payload []byte) error {
	buffer := bytes.NewBuffer(payload)
	messageType, err := buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read message type: %w", err)
	}
	switch messageType {
	case REALTIME_CAR_UPDATE:
		updateRealtimeCarUpdate(buffer, t.realtimeCarUpdate)
	case REALTIME_UPDATE:
		// silently drop this one
		break
	default:
		return fmt.Errorf("we don't currently support this type of message: '%d' => %+v", messageType, payload)
	}

	return nil
}
