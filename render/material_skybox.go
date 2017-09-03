package render

func newSkyboxMaterial() *Material {
	return newMaterial(skyboxVertexShader, skyboxFragmentShader)
}

const skyboxVertexShader string = `#version 120

uniform mat4 projectionMatrixIn;
uniform mat4 viewMatrixIn;

attribute vec3 coordIn;
attribute vec2 texCoordIn;

varying vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	// assure that translations are ignored by setting w to 0.0
	vec4 viewPosition = viewMatrixIn * vec4(coordIn, 0.0);
	// restore w to 1.0 so that projection works
	gl_Position = projectionMatrixIn * vec4(viewPosition.xyz, 1.0);
}
`

const skyboxFragmentShader string = `#version 120

uniform sampler2D diffuseTextureIn;

varying vec2 texCoordInOut;

void main()
{
	gl_FragColor = texture2D(diffuseTextureIn, texCoordInOut);
}
`
