package asset

import (
	"encoding/gob"
	"fmt"
	"io"
)

type TextureFormat uint8

const (
	TextureFormatRGBA TextureFormat = iota
	TextureFormatBGRA
	TextureFormatRGB
)

type TextureSide int

const (
	TextureSideFront TextureSide = iota
	TextureSideBack
	TextureSideLeft
	TextureSideRight
	TextureSideTop
	TextureSideBottom
)

type TwoDTexture struct {
	Width  uint16
	Height uint16
	Format TextureFormat
	Data   []byte
}

func NewTwoDTextureDecoder() *TwoDTextureDecoder {
	return &TwoDTextureDecoder{}
}

type TwoDTextureDecoder struct{}

func (d *TwoDTextureDecoder) Decode(in io.Reader) (*TwoDTexture, error) {
	var texture TwoDTexture
	if err := gob.NewDecoder(in).Decode(&texture); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &texture, nil
}

func NewTwoDTextureEncoder() *TwoDTextureEncoder {
	return &TwoDTextureEncoder{}
}

type TwoDTextureEncoder struct{}

func (e *TwoDTextureEncoder) Encode(out io.Writer, texture *TwoDTexture) error {
	if err := gob.NewEncoder(out).Encode(texture); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}

type CubeTexture struct {
	Dimension uint16
	Format    TextureFormat
	Sides     [6]CubeTextureSide
}

type CubeTextureSide struct {
	Data []byte
}

func NewCubeTextureDecoder() *CubeTextureDecoder {
	return &CubeTextureDecoder{}
}

type CubeTextureDecoder struct{}

func (d *CubeTextureDecoder) Decode(in io.Reader) (*CubeTexture, error) {
	var texture CubeTexture
	if err := gob.NewDecoder(in).Decode(&texture); err != nil {
		return nil, fmt.Errorf("failed to decode gob stream: %w", err)
	}
	return &texture, nil
}

func NewCubeTextureEncoder() *CubeTextureEncoder {
	return &CubeTextureEncoder{}
}

type CubeTextureEncoder struct{}

func (e *CubeTextureEncoder) Encode(out io.Writer, texture *CubeTexture) error {
	if err := gob.NewEncoder(out).Encode(texture); err != nil {
		return fmt.Errorf("failed to encode gob stream: %w", err)
	}
	return nil
}
