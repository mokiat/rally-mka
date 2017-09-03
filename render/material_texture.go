package render

func newTextureMaterial() *Material {
	return newMaterial(textureVertexShader, textureFragmentShader)
}

const textureVertexShader string = `#version 120

uniform mat4 projectionMatrixIn;
uniform mat4 modelMatrixIn;
uniform mat4 viewMatrixIn;

attribute vec3 coordIn;
attribute vec2 texCoordIn;

varying vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * (modelMatrixIn * vec4(coordIn, 1.0)));
}
`

const textureFragmentShader string = `#version 120

uniform sampler2D diffuseTextureIn;

varying vec2 texCoordInOut;

void main()
{
	gl_FragColor = texture2D(diffuseTextureIn, texCoordInOut);
}
`
