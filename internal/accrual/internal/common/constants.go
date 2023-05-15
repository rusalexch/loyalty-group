package common

const (
	Registered = "REGISTERED"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)

type RewardType string

const (
	Percentage RewardType = "%"
	Fixed RewardType = "pt"
)
