package rclone

import "github.com/dustin/go-humanize"

type ProgressMessage struct {
	Level  string               `json:"level"`
	Msg    string               `json:"msg"`
	Source string               `json:"source"`
	Stats  ProgressMessageStats `json:"stats"`
}

type ProgressMessageStats struct {
	Bytes          int64                              `json:"bytes"`
	Checks         int64                              `json:"checks"`
	DeletedDirs    int64                              `json:"deletedDirs"`
	Deletes        int64                              `json:"deletes"`
	ElapsedTime    float64                            `json:"elapsedTime"`
	Errors         int64                              `json:"errors"`
	Eta            int64                              `json:"eta"`
	FatalError     bool                               `json:"fatalError"`
	Renames        int64                              `json:"renames"`
	RetryError     bool                               `json:"retryError"`
	Speed          float64                            `json:"speed"`
	TotalBytes     int64                              `json:"totalBytes"`
	TotalChecks    int64                              `json:"totalChecks"`
	TotalTransfers int64                              `json:"totalTransfers"`
	TransferTime   int64                              `json:"transferTime"`
	Transferring   []ProgressMessageStatsTransferring `json:"transferring"`
	Transfers      int64                              `json:"transfers"`
}

type ProgressMessageStatsTransferring struct {
	Bytes      int64   `json:"bytes"`
	Eta        int64   `json:"eta"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Percentage int     `json:"percentage"`
	Size       int64   `json:"size"`
	Speed      float64 `json:"speed"`
	SpeedAvg   float64 `json:"speedAvg"`
}

func (pm *ProgressMessage) IsValid() bool {
	return pm.Stats.Bytes > 0 && pm.Stats.TotalBytes > 0
}

// Progress of whole transfer, float 0 to 1
func (pm *ProgressMessage) Progress() float64 {
	return float64(pm.Stats.Bytes) / float64(pm.Stats.TotalBytes)
}

func (pm *ProgressMessage) HumanSpeed() string {
	return humanize.IBytes(uint64(pm.Stats.Speed)) + "/s"
}
