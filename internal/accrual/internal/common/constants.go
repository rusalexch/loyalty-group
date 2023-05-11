package common

type OrderStatus string

const (
	Registered OrderStatus = "REGISTERED"
	Invalid    OrderStatus = "INVALID"
	Processing OrderStatus = "PROCESSING"
	Processed  OrderStatus = "PROCESSED"
)

type RewardType string

const (
	Percentage RewardType = "%"
	Fixed RewardType = "pt"
)
