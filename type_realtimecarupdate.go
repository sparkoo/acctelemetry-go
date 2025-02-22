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

func updateRealtimeCarUpdate(payload *bytes.Buffer, toUpdate *RealtimeCarUpdate) {
	toUpdate.CarIndex, _ = readUint16(payload)
	toUpdate.DriverIndex, _ = readUint16(payload)
	toUpdate.DriverCount, _ = payload.ReadByte()
	toUpdate.Gear, _ = payload.ReadByte()
	toUpdate.WorldPosX, _ = readFloat32(payload)
	toUpdate.WorldPosY, _ = readFloat32(payload)
	toUpdate.Yaw, _ = readFloat32(payload)
	toUpdate.CarLocation, _ = payload.ReadByte()
	toUpdate.Kmh, _ = readUint16(payload)
	toUpdate.Position, _ = readUint16(payload)
	toUpdate.CupPosition, _ = readUint16(payload)
	toUpdate.TrackPosition, _ = readUint16(payload)
	toUpdate.TrackRelativePosition, _ = readFloat32(payload)
	toUpdate.Laps, _ = readUint16(payload)
	toUpdate.Delta, _ = readInt32(payload)
	updateLap(payload, toUpdate.BestSessionLap)
	updateLap(payload, toUpdate.LastLap)
	updateLap(payload, toUpdate.CurrentLap)
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
