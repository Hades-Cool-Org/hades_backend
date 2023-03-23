package model

type Product struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Details       string `json:"details"`
	Image         string `json:"image_url"`
	MeasuringUnit string `json:"measuring_unit"`
}
