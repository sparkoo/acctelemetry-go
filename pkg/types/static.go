package types

type AccStatic struct {
	SmVersion              [15]rune
	AcVersion              [15]rune
	NumberOfSessions       int32
	NumCars                int32
	CarModel               [33]rune
	Track                  [33]rune
	PlayerName             [33]rune
	PlayerSurname          [33]rune
	PlayerNick             [33]rune
	SectorCount            int32
	MaxTorque              float32
	MaxPower               float32
	MaxRpm                 int32
	MaxFuel                float32
	SuspensionMaxTravel    [4]float32
	TyreRadius             [4]float32
	MaxTurboBoost          [4]float32
	Deprecated1            float32
	Deprecated2            float32
	PenaltiesEnabled       int32
	AidFuelRate            float32
	AidTireRate            float32
	AidMechanicalDamage    float32
	AidAllowTyreBlankets   int32
	AidStability           float32
	AidAutoClutch          int32
	AidAutoBlip            int32
	HasDRS                 int32
	HasERS                 int32
	HasKERS                int32
	KersMaxJ               float32
	EngineBrakeSettingsCnt int32
	ErsPowerControllerCnt  int32
	TrackSplineLength      float32
	TrackConfiguration     [33]rune
	ErsMaxJ                float32
	IsTimedRace            int32
	HasExtraLap            int32
	CarSkin                [33]rune
	ReversedGridPositions  int32
	PitWindowStart         int32
	PitWindowEnd           int32
	IsOnline               int32
}
