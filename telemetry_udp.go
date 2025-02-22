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
	n, remoteAddr, err := telemetry.udpConnection.ReadFromUDP(inBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("UDP read timeout, ACC may not be running: %w", netErr)
		} else {
			return fmt.Errorf("UDP read failed: %w", err)
		}
	}
	fmt.Printf("Connected to ACC, bytes '%d', addr '%+v': '%+v'", n, remoteAddr, inBuffer)
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
