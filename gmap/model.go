package gmap

type Coordinate struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type geocodingResult struct {
	Results []struct {
		Geometry struct {
			Coordinate Coordinate `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}
