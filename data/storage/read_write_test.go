package storage_test

import (
	"bytes"
	"strings"

	. "github.com/mokiat/rally-mka/data/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {

	var buffer *bytes.Buffer
	var writer TypedWriter
	var reader TypedReader

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		writer = NewTypedWriter(buffer)
		reader = NewTypedReader(buffer)
	})

	Context("when byte array is written", func() {
		var writtenValue []byte

		BeforeEach(func() {
			writtenValue = []byte("hello")
			err := writer.WriteBytes(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("correct number of bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(len(writtenValue)))
		})

		It("is possible to read the value back", func() {
			readValue := make([]byte, len(writtenValue))
			err := reader.ReadBytes(readValue)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when byte is written", func() {
		var writtenValue byte

		BeforeEach(func() {
			writtenValue = 0x3C
			err := writer.WriteByte(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("1 byte is written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(1))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadByte()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when uint8 is written", func() {
		var writtenValue uint8

		BeforeEach(func() {
			writtenValue = 0x2D
			err := writer.WriteUInt8(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("1 byte is written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(1))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadUInt8()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when int8 is written", func() {
		var writtenValue int8

		BeforeEach(func() {
			writtenValue = -0x2D
			err := writer.WriteInt8(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("1 byte is written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(1))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadInt8()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when uint16 is written", func() {
		var writtenValue uint16

		BeforeEach(func() {
			writtenValue = 0xFA1F
			err := writer.WriteUInt16(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("2 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(2))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadUInt16()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when int16 is written", func() {
		var writtenValue int16

		BeforeEach(func() {
			writtenValue = -0x5A1F
			err := writer.WriteInt16(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("2 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(2))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadInt16()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when uint32 is written", func() {
		var writtenValue uint32

		BeforeEach(func() {
			writtenValue = 0xC5B8F2E4
			err := writer.WriteUInt32(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("4 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(4))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadUInt32()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when int32 is written", func() {
		var writtenValue int32

		BeforeEach(func() {
			writtenValue = -0x15B8F2E4
			err := writer.WriteInt32(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("4 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(4))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadInt32()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when uint64 is written", func() {
		var writtenValue uint64

		BeforeEach(func() {
			writtenValue = 0xC5B8F2E4A2B3C4D5
			err := writer.WriteUInt64(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("8 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(8))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadUInt64()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when int64 is written", func() {
		var writtenValue int64

		BeforeEach(func() {
			writtenValue = -0x25B8F2E4A2B3C4D5
			err := writer.WriteInt64(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("8 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(8))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadInt64()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when float32 is written", func() {
		var writtenValue float32

		BeforeEach(func() {
			writtenValue = 1.548952
			err := writer.WriteFloat32(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("4 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(4))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadFloat32()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when float64 is written", func() {
		var writtenValue float64

		BeforeEach(func() {
			writtenValue = 8.12237651
			err := writer.WriteFloat64(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("8 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(8))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadFloat64()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	produceString := func(length int) string {
		return strings.Repeat("a", length)
	}

	Context("when 8 bit long string is written", func() {
		var writtenValue string

		BeforeEach(func() {
			writtenValue = produceString(255)
			err := writer.WriteString8(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("correct number of bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(1 + 255))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadString8()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when 16 bit long string is written", func() {
		var writtenValue string

		BeforeEach(func() {
			writtenValue = produceString(65535)
			err := writer.WriteString16(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("2+255 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(2 + 65535))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadString16()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

	Context("when 32 bit long string is written", func() {
		var writtenValue string

		BeforeEach(func() {
			writtenValue = produceString(65536)
			err := writer.WriteString32(writtenValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("2+255 bytes are written", func() {
			Ω(buffer.Bytes()).Should(HaveLen(4 + 65536))
		})

		It("is possible to read the value back", func() {
			readValue, err := reader.ReadString32()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readValue).Should(Equal(writtenValue))
		})
	})

})
