package graphics

type Program struct {
	ID                       uint32
	ProjectionMatrixLocation int32
	ViewMatrixLocation       int32
	ModelMatrixLocation      int32
	DiffuseTextureLocation   int32
	SkyboxTextureLocation    int32
}
