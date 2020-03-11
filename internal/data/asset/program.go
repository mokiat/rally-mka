package asset

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Program struct {
	VertexSourceCode   string
	FragmentSourceCode string
}

func NewProgramDecoder() *ProgramDecoder {
	return &ProgramDecoder{}
}

type ProgramDecoder struct{}

func (d *ProgramDecoder) Decode(in io.Reader) (*Program, error) {
	var program Program
	if err := gob.NewDecoder(in).Decode(&program); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &program, nil
}

func NewProgramEncoder() *ProgramEncoder {
	return &ProgramEncoder{}
}

type ProgramEncoder struct{}

func (e *ProgramEncoder) Encode(out io.Writer, program *Program) error {
	if err := gob.NewEncoder(out).Encode(program); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}
