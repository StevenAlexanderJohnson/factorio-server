package models

type AutoplaceControlSetting struct {
	Frequency float64 `json:"frequency,omitempty"`
	Size      float64 `json:"size,omitempty"`
	Richness  float64 `json:"richness,omitempty"`
}

type CliffSettings struct {
	Name                   string  `json:"name,omitempty"`
	CliffElevation0        float64 `json:"cliff_elevation_0,omitempty"`
	CliffElevationInterval float64 `json:"cliff_elevation_interval,omitempty"`
	Richness               float64 `json:"richness,omitempty"`
}

type StartingPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MapGenSettings struct {
	Width                   int                                `json:"width"`
	Height                  int                                `json:"height"`
	StartingArea            float64                            `json:"starting_area"`
	PeacefulMode            bool                               `json:"peaceful_mode"`
	AutoplaceControls       map[string]AutoplaceControlSetting `json:"autoplace_controls,omitempty"`
	CliffSettings           CliffSettings                      `json:"cliff_settings,omitempty"`
	PropertyExpressionNames map[string]string                  `json:"property_expression_names,omitempty"`
	StartingPoints          []StartingPoint                    `json:"starting_points,omitempty"`
	Seed                    *uint32                            `json:"seed"`
}
