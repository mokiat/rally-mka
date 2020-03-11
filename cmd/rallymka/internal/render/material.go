package render

import "github.com/mokiat/go-whiskey-gl/shader"

type Material struct {
	program            *shader.Program
	vertexSourceCode   string
	fragmentSourceCode string

	coordLocation    uint32
	normalLocation   uint32
	colorLocation    uint32
	texCoordLocation uint32

	projectionMatrixLocation int32
	modelMatrixLocation      int32
	viewMatrixLocation       int32
	diffuseTextureLocation   int32
	skyboxTextureLocation    int32
}

func newMaterial(vertexSrc, fragmentSrc string) *Material {
	return &Material{
		vertexSourceCode:   vertexSrc,
		fragmentSourceCode: fragmentSrc,
	}
}

func (m *Material) Generate() error {
	vertexShader, err := initVertexShader(m.vertexSourceCode)
	if err != nil {
		return err
	}
	fragmentShader, err := initFragmentShader(m.fragmentSourceCode)
	if err != nil {
		return err
	}
	m.program = shader.NewProgram()
	if err := m.program.Allocate(); err != nil {
		return err
	}
	m.program.AttachVertexShader(vertexShader)
	m.program.AttachFragmentShader(fragmentShader)
	if err := m.program.LinkProgram(); err != nil {
		return err
	}

	m.coordLocation = m.program.GetAttributeLocation("coordIn")
	m.normalLocation = m.program.GetAttributeLocation("normalIn")
	m.colorLocation = m.program.GetAttributeLocation("colorIn")
	m.texCoordLocation = m.program.GetAttributeLocation("texCoordIn")

	m.projectionMatrixLocation = m.program.GetUniformLocation("projectionMatrixIn")
	m.modelMatrixLocation = m.program.GetUniformLocation("modelMatrixIn")
	m.viewMatrixLocation = m.program.GetUniformLocation("viewMatrixIn")
	m.diffuseTextureLocation = m.program.GetUniformLocation("diffuseTextureIn")
	m.skyboxTextureLocation = m.program.GetUniformLocation("skyboxTextureIn")
	return nil
}

func initVertexShader(source string) (*shader.VertexShader, error) {
	vertexShader := shader.NewVertexShader()
	if err := vertexShader.Allocate(); err != nil {
		return nil, err
	}
	vertexShader.SetSourceCode(source)
	if err := vertexShader.Compile(); err != nil {
		return nil, err
	}
	return vertexShader, nil
}

func initFragmentShader(source string) (*shader.FragmentShader, error) {
	fragmentShader := shader.NewFragmentShader()
	if err := fragmentShader.Allocate(); err != nil {
		return nil, err
	}
	fragmentShader.SetSourceCode(source)
	if err := fragmentShader.Compile(); err != nil {
		return nil, err
	}
	return fragmentShader, nil
}
