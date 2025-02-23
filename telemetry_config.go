package acctelemetry

type accTelemetryConfig struct {
	EnableUdp                   bool
	UdpIpPort                   string
	UdpConnectionPassword       string
	udpCommandPassword          string
	udpDisplayName              string
	udpRealtimeUpdateIntervalMS int32
}

func DefaultConfig() *accTelemetryConfig {
	return &accTelemetryConfig{
		EnableUdp: false,
	}
}

func DefaultUdpConfig() *accTelemetryConfig {
	return UdpConfig("127.0.0.1:9000", "asd")
}

func UdpConfig(ipPort string, password string) *accTelemetryConfig {
	return &accTelemetryConfig{
		EnableUdp:                   true,
		UdpIpPort:                   ipPort,
		UdpConnectionPassword:       password,
		udpCommandPassword:          "",
		udpDisplayName:              "RaceMate",
		udpRealtimeUpdateIntervalMS: 100,
	}
}
