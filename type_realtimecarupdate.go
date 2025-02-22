package acctelemetry

import (
	"bytes"
)

type RealtimeCarUpdate struct {
	CarIndex              uint16
	DriverIndex           uint16
	DriverCount           byte
	Gear                  byte
	WorldPosX             float32
	WorldPosY             float32
	Yaw                   float32
	CarLocation           byte // 0=NONE, 1=TRACK, 2=PITLANE, 3=PITENTRY, 4=PITEXIT
	Kmh                   uint16
	Position              uint16
	CupPosition           uint16
	TrackPosition         uint16
	TrackRelativePosition float32
	Laps                  uint16
	Delta                 int32
	BestSessionLap        *LapInfo
	LastLap               *LapInfo
	CurrentLap            *LapInfo
}

type LapInfo struct {
	LaptimeMs      int32
	CarIndex       uint16
	DriverIndex    uint16
	SplitCount     byte
	Splits         [8]int32 // let's assume no track will have more than 8 splits (all should have 3...)
	IsInvalid      byte
	InValidForBest byte
	IsOutlap       byte
	IsInlaop       byte
}

func updateRealtimeCarUpdate(payload *bytes.Buffer) *RealtimeCarUpdate {
	carUpdate := &RealtimeCarUpdate{
		BestSessionLap: &LapInfo{Splits: [8]int32{}},
		LastLap:        &LapInfo{Splits: [8]int32{}},
		CurrentLap:     &LapInfo{Splits: [8]int32{}},
	}
	carUpdate.CarIndex, _ = readUint16(payload)
	carUpdate.DriverIndex, _ = readUint16(payload)
	carUpdate.DriverCount, _ = payload.ReadByte()
	carUpdate.Gear, _ = payload.ReadByte()
	carUpdate.WorldPosX, _ = readFloat32(payload)
	carUpdate.WorldPosY, _ = readFloat32(payload)
	carUpdate.Yaw, _ = readFloat32(payload)
	carUpdate.CarLocation, _ = payload.ReadByte()
	carUpdate.Kmh, _ = readUint16(payload)
	carUpdate.Position, _ = readUint16(payload)
	carUpdate.CupPosition, _ = readUint16(payload)
	carUpdate.TrackPosition, _ = readUint16(payload)
	carUpdate.TrackRelativePosition, _ = readFloat32(payload)
	carUpdate.Laps, _ = readUint16(payload)
	carUpdate.Delta, _ = readInt32(payload)
	updateLap(payload, carUpdate.BestSessionLap)
	updateLap(payload, carUpdate.LastLap)
	updateLap(payload, carUpdate.CurrentLap)

	return carUpdate
}

func updateLap(payload *bytes.Buffer, lap *LapInfo) {
	lap.LaptimeMs, _ = readInt32(payload)
	lap.CarIndex, _ = readUint16(payload)
	lap.DriverIndex, _ = readUint16(payload)
	lap.SplitCount, _ = payload.ReadByte()
	for i := 0; i < int(lap.SplitCount); i++ {
		lap.Splits[i], _ = readInt32(payload)
	}
	lap.IsInvalid, _ = payload.ReadByte()
	lap.InValidForBest, _ = payload.ReadByte()
	lap.IsOutlap, _ = payload.ReadByte()
	lap.IsInlaop, _ = payload.ReadByte()
}
