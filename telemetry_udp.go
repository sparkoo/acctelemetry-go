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
	REQUEST_ENTRY_LIST             byte = 10
	REQUEST_TRACK_DATA             byte = 11
	CHANGE_HUD_PAGE                byte = 49
	CHANGE_FOCUS                   byte = 50
	INSTANT_REPLAY_REQUEST         byte = 51
	PLAY_MANUAL_REPLAY_HIGHLIGHT   byte = 52
	SAVE_MANUAL_REPLAY_HIGHLIGHT   byte = 60
)

const (
	REGISTRATION_RESULT byte = 1
	REALTIME_UPDATE     byte = 2
	REALTIME_CAR_UPDATE byte = 3
	ENTRY_LIST          byte = 4
	ENTRY_LIST_CAR      byte = 6
	TRACK_DATA          byte = 5
	BROADCASTING_EVENT  byte = 7
)

type connectionResult struct {
	connectionId      int32
	connectionSuccess bool
	readOnly          bool
	errorMessage      string
}

func (telemetry *accTelemetry) connect() error {
	outBuffer := bytes.NewBuffer([]byte{})
	var writeErr error
	writeErr = outBuffer.WriteByte(REGISTER_COMMAND_APPLICATION)
	writeErr = outBuffer.WriteByte(BROADCASTING_PROTOCOL_VERSION)
	writeErr = writeString(outBuffer, telemetry.config.UdpDisplayName)
	writeErr = writeString(outBuffer, telemetry.config.UdpConnectionPassword)
	writeErr = binary.Write(outBuffer, binary.LittleEndian, telemetry.config.UdpRealtimeUpdateIntervalMS)
	writeErr = writeString(outBuffer, telemetry.config.UdpCommandPassword)

	if writeErr != nil {
		return fmt.Errorf("failed to write connection data to byte buffer: %w", writeErr)
	}

	_, sendErr := telemetry.udpConnection.Write(outBuffer.Bytes())
	if sendErr != nil {
		return fmt.Errorf("failed to send connection message: %w", sendErr)
	}
	telemetry.udpConnection.SetReadDeadline(time.Now().Add(1 * time.Second))
	inBuffer := make([]byte, 1024)
	n, _, err := telemetry.udpConnection.ReadFromUDP(inBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("UDP read timeout, ACC may not be running: %w", netErr)
		} else {
			return fmt.Errorf("UDP read failed: %w", err)
		}
	}
	connectionResult, err := readConnectionResult(n, bytes.NewBuffer(inBuffer))
	if err != nil {
		return fmt.Errorf("failed to read connection response: %w", err)
	}
	fmt.Printf("Connected to ACC: '%+v'", connectionResult)
	return nil
}

func writeString(buffer *bytes.Buffer, str string) error {
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, uint16(len(str)))

	var writeErr error
	_, writeErr = buffer.Write(lengthBytes)
	_, writeErr = buffer.Write([]byte(str))
	if writeErr != nil {
		return fmt.Errorf("failed to write string '%s' to byte buffer: %w", str, writeErr)
	}
	return nil
}

// response format:
// 1 byte - message type
// 4 bytes - int32 connectionId
// 1 byte - connectionSuccess `byte > 0`
// 1 byte - readonly `byte == 0`
// 2 bytes - error length => N
// N bytes - error message
func readConnectionResult(length int, payload *bytes.Buffer) (*connectionResult, error) {
	result := &connectionResult{}

	// read type
	messageType, err := payload.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read response type: %w", err)
	} else if messageType != REGISTRATION_RESULT {
		return nil, fmt.Errorf("should've received REGISTRATION_RESULT['1'] message type, but got '%d'", messageType)
	}

	// read connectionId
	connectionIdBytes := make([]byte, 4)
	_, err = payload.Read(connectionIdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read connectionId: %w", err)
	}

	var connectionId int32
	err = binary.Read(bytes.NewReader(connectionIdBytes), binary.LittleEndian, &connectionId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert connectionId: %w", err)
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
	errorLengthBytes := make([]byte, 2)
	_, err = payload.Read(errorLengthBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read error length: %w", err)
	}
	var errorLength uint16
	err = binary.Read(bytes.NewReader(errorLengthBytes), binary.LittleEndian, &errorLength)
	if err != nil {
		return nil, fmt.Errorf("failed to convert errorLength: %w", err)
	}

	if errorLength > 0 {
		result.errorMessage = "some error"
	}

	return result, nil
}
