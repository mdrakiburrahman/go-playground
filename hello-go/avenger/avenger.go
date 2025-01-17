package avenger

// Avenger represents a single hero
type Avenger struct {
	RealName string `json:"real_name"`
	HeroName string `json:"hero_name"`
	Planet   string `json:"planet"`
	Alive    bool   `json:"alive"`
}

// IsAlive makes an Avenger alive
func (a *Avenger) IsAlive() {
	a.Alive = true
}
