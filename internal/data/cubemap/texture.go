package cubemap

const Version uint8 = 1

type Texture struct {
	Format    DataFormat
	Dimension uint16
	Sides     [6]TextureSide
}

type TextureSide struct {
	Data []byte
}

type Side int

const (
	SideFront Side = iota
	SideBack
	SideLeft
	SideRight
	SideTop
	SideBottom
)

type DataFormat uint8

const (
	DataFormatRGBA DataFormat = iota
	DataFormatBGRA
)
