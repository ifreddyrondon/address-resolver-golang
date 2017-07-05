package addresses

type Address struct {
	Address string  `json:"address"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

type Addresses []Address
