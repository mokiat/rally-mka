package resource

import (
	"fmt"
	"io"
	"io/ioutil"
)

type Shader struct {
	SourceCode string
}

func NewShaderDecoder() *ShaderDecoder {
	return &ShaderDecoder{}
}

type ShaderDecoder struct{}

func (d *ShaderDecoder) Decode(in io.Reader) (*Shader, error) {
	contents, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read all data: %w", err)
	}
	return &Shader{
		SourceCode: string(contents),
	}, nil
}
