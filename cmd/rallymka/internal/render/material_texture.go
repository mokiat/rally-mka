package render

func newTextureMaterial() *Material {
	return newMaterial(textureVertexShader, textureFragmentShader)
}

const textureVertexShader string = `#version 410

uniform mat4 projectionMatrixIn;
uniform mat4 modelMatrixIn;
uniform mat4 viewMatrixIn;

in vec3 coordIn;
in vec2 texCoordIn;

smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * (modelMatrixIn * vec4(coordIn, 1.0)));
}
`

const textureFragmentShader string = `#version 410

uniform sampler2D diffuseTextureIn;

smooth in vec2 texCoordInOut;
layout(location = 0) out vec4 fragmentColor;

void main()
{
	vec4 color = texture(diffuseTextureIn, texCoordInOut);
	if (color.a < 0.9) {
		discard;
	}
	fragmentColor = color;
}
`
