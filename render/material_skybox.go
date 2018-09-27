package render

func newSkyboxMaterial() *Material {
	return newMaterial(skyboxVertexShader, skyboxFragmentShader)
}

const skyboxVertexShader string = `#version 410

uniform mat4 projectionMatrixIn;
uniform mat4 viewMatrixIn;

in vec3 coordIn;
in vec2 texCoordIn;

smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	// assure that translations are ignored by setting w to 0.0
	vec4 viewPosition = viewMatrixIn * vec4(coordIn, 0.0);
	// restore w to 1.0 so that projection works
	gl_Position = projectionMatrixIn * vec4(viewPosition.xyz, 1.0);
}
`

const skyboxFragmentShader string = `#version 410

uniform sampler2D diffuseTextureIn;

smooth in vec2 texCoordInOut;
layout(location = 0) out vec4 fragmentColor;

void main()
{
	fragmentColor = texture(diffuseTextureIn, texCoordInOut);
}
`
