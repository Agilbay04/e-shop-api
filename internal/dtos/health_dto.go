package dto

type HealthResponse struct {
	Status   string         `json:"status"`
	Service  string         `json:"service"`
	Version  string         `json:"version"`
	Details  map[string]any `json:"details,omitempty"`
}

type ReadinessCheck struct {
	Status  string         `json:"status"`
	Checks  map[string]any `json:"checks"`
}

type ComponentStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}