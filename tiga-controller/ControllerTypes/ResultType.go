package ControllerTypes

type RemoteConfig struct {
	Enable    bool                `json:"enable"`
	Url       map[string]string   `json:"url"`
	Metrics   []string            `json:"metrics"`
	Blacklist []map[string]string `json:"blacklist"`
}

type LocalMetrics struct {
	Enable  bool     `json:"enable"`
	Metrics []string `json:"metrics"`
}

type LocalConfig struct {
	Container LocalMetrics `json:"container"`
	Node      LocalMetrics `json:"node"`
}

type RespConfig struct {
	Local  LocalConfig  `json:"local"`
	Remote RemoteConfig `json:"remote"`
}
