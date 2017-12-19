package render

func newSkycubeMaterial() *Material {
	return newMaterial(skycubeVertexShader, skycubeFragmentShader)
}

const skycubeVertexShader string = `#version 120

uniform mat4 projectionMatrixIn;
uniform mat4 viewMatrixIn;

attribute vec3 coordIn;

varying vec3 texCoordInOut;

void main()
{
	// we optimize by using vertex coords as cube texture coords
	// we need to normalize the coords so that clamping works correctly
	// additionally, we need to flip the coords. opengl is broken
	// in that it uses renderman coordinate system for cube maps, 
	// contrary to the rest of the opengl api. epic fail!
	texCoordInOut = -normalize(coordIn);

	// assure that translations are ignored by setting w to 0.0
	vec4 viewPosition = viewMatrixIn * vec4(coordIn, 0.0);

	// restore w to 1.0 so that projection works
	gl_Position = projectionMatrixIn * vec4(viewPosition.xyz, 1.0);
}
`

const skycubeFragmentShader string = `#version 120

uniform samplerCube skycubeTextureIn;

varying vec3 texCoordInOut;

void main()
{
	// gl_FragColor = textureCube(skycubeTextureIn, texCoordInOut);
	gl_FragColor = vec4(textureCube(skycubeTextureIn, texCoordInOut).rgb, 1.0);
	// gl_FragColor = vec4(texCoordInOut, 1.0);
}
`
