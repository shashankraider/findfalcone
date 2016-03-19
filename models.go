package main

type City struct {
	Name     string `json:"name"`
	Distance int    `json:"distance"`
}

type Cities []City

type TransportVehicle struct {
	Name        string `json:"name"`
	TotalNo     int    `json:"total_no"`
	MaxDistance int    `json:"max_distance"`
	Speed       int    `json:"speed"`
}

type TransportVehicles []TransportVehicle

type FindFalconeReq struct {
	Token        string   `json:"token"`
	CityNames  []string `json:"city_names"`
	TransportVehicleNames []string `json:"transportvehicle_names"`
}
