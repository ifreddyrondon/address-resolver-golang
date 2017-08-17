package gmap

type Coordinate struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Geometry struct {
	Coordinate Coordinate `json:"location"`
}

type geocodingResult struct {
	Results []struct {
		Geometry Geometry `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}
