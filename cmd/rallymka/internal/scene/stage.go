package scene

type Stage struct {
	Sky *Skybox
}

func NewStage() *Stage {
	return &Stage{}
}

func (s *Stage) GetSky(camera *Camera) (*Skybox, bool) {
	// TODO: this implementation should iterate the BSP tree / Octree
	// to find which skybox is relevant for the current camera position
	// IDEA: Does blending between skyboxes at scene region boundaries make sense?
	return s.Sky, s.Sky != nil
}
