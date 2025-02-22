package acctelemetry

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

func readInt32(buffer *bytes.Buffer) (int32, error) {
	readBytes := make([]byte, 4)
	_, err := buffer.Read(readBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to read int32: %w", err)
	}

	var result int32
	err = binary.Read(bytes.NewReader(readBytes), binary.LittleEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to convert int32: %w", err)
	}
	return result, nil
}

func readUint16(buffer *bytes.Buffer) (uint16, error) {
	readBytes := make([]byte, 2)
	_, err := buffer.Read(readBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to read uint: %w", err)
	}
	var result uint16
	err = binary.Read(bytes.NewReader(readBytes), binary.LittleEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to convert uint: %w", err)
	}

	return result, nil
}

func readFloat32(buffer *bytes.Buffer) (float32, error) {
	readBytes := make([]byte, 4)
	_, err := buffer.Read(readBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to read float32: %w", err)
	}

	var result float32
	err = binary.Read(bytes.NewReader(readBytes), binary.LittleEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to convert float32: %w", err)
	}
	return result, nil
}
