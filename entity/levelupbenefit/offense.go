package levelupbenefit

// Offense describes all of the offensive benefits of leveling up.
type Offense struct {
	Aim      int `json:"aim" yaml:"aim"`
	Strength int `json:"strength" yaml:"strength"`
	Mind     int `json:"mind" yaml:"mind"`
}
