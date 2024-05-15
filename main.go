package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// OmcTextDecoder struct
type OmcTextDecoder struct {
	XMLHeader  string
	SaltLength int
	Salts      []int8
	Shifts     []byte
}

// NewOmcTextDecoder constructor
func NewOmcTextDecoder() *OmcTextDecoder {
	return &OmcTextDecoder{
		XMLHeader:  "<?xml",
		SaltLength: 256,
		Salts: []int8{65, -59, 33, -34, 107, 28, -107, 55, 78, 17, -81, 6, -80, -121, -35, -23, 72, 122, -63, -43,
			68, 119, -78, -111, -60, 31, 60, 57, 92, -88, -100, -69, -106, 91, 69, 93, 110, 23, 93, 53, -44, -51,
			64, -80, 46, 2, -4, 12, -45, 80, -44, -35, -111, -28, -66, -116, 39, 2, -27, -45, -52, 125, 39, 66, -90,
			63, -105, -67, 84, -57, -4, -4, 101, -90, 81, 10, -33, 1, 67, -57, -71, 18, -74, 102, 96, -89, 64, -17,
			54, -94, -84, -66, 14, 119, 121, 2, -78, -79, 89, 63, 93, 109, -78, -51, 66, -36, 32, 86, 3, -58, -15,
			92, 58, 2, -89, -80, -13, -1, 122, -4, 48, 63, -44, 59, 100, -42, -45, 59, -7, -17, -54, 34, -54, 71,
			-64, -26, -87, -80, -17, -44, -38, -112, 70, 10, -106, 95, -24, -4, -118, 45, -85, -13, 85, 25, -102,
			-119, 13, -37, 116, 46, -69, 59, 42, -90, -38, -105, 101, -119, -36, 97, -3, -62, -91, -97, -125, 17, 14,
			106, -72, -119, 99, 111, 20, 18, -27, 113, 64, -24, 74, -60, -100, 26, 56, -44, -70, 12, -51, -100, -32,
			-11, 26, 48, -117, 98, -93, 51, -25, -79, -31, 97, 87, -105, -64, 7, -13, -101, 33, -122, 5, -104, 89,
			-44, -117, 63, -80, -6, -71, -110, -29, -105, 116, 107, -93, 91, -41, -13, 20, -115, -78, 43, 79, -122,
			6, 102, -32, 52, -118, -51, 72, -104, 41, -38, 124, 72, -126, -35},
		Shifts: []byte{1, 1, 0, 2, 2, 4, 5, 0, 4, 7, 1, 6, 5, 3, 3, 1, 2, 5, 0, 6, 2, 2, 4, 2, 2, 3, 0, 2, 1, 2, 4, 3, 4, 0,
			0, 0, 3, 5, 3, 1, 6, 5, 6, 1, 1, 1, 0, 0, 3, 2, 7, 7, 5, 6, 7, 3, 5, 1, 0, 7, 6, 3, 6, 5, 4, 5, 3, 5, 1, 3, 3,
			1, 5, 4, 1, 0, 0, 2, 6, 6, 6, 6, 4, 0, 1, 1, 0, 5, 5, 4, 2, 4, 6, 1, 7, 1, 2, 1, 1, 6, 5, 4, 7, 6, 5, 1, 6, 7,
			0, 2, 6, 3, 1, 7, 1, 1, 7, 4, 0, 4, 2, 5, 3, 1, 1, 5, 6, 0, 3, 5, 3, 6, 5, 7, 2, 5, 6, 6, 2, 2, 3, 6, 0, 4, 3,
			2, 0, 2, 2, 3, 5, 3, 3, 2, 5, 5, 5, 1, 3, 1, 1, 1, 4, 5, 1, 6, 2, 4, 7, 1, 4, 6, 0, 6, 4, 3, 2, 6, 1, 6, 3, 2, 1, 6, 7, 3, 2, 1, 1, 5, 6, 7, 2, 2, 2, 7, 4, 6, 7, 5, 3, 1, 4, 2, 7, 1, 6, 2, 4, 1, 5, 6,
			5, 4, 5, 0, 1, 1, 6, 3, 7, 2, 0, 2, 5, 0, 1, 3, 3, 2, 6, 7, 7, 2, 5, 6, 0, 4, 1, 2, 5, 3, 7, 6, 5, 2, 5, 2,
			0, 1, 3, 1, 4, 3, 4, 2},
	}
}

// _decode method
func (decoder *OmcTextDecoder) _decode(bArr []byte) []byte {
	bArr2 := make([]byte, len(bArr))
	for i := range bArr {
		b := bArr[i]
		i2 := i % 256
		b2 := decoder.Shifts[i2]
		b3 := (b >> (8 - b2)) | (b << b2)
		bArr2[i] = b3 ^ byte(decoder.Salts[i2])
	}
	return bArr2
}

// _decompressGzip method
func (decoder *OmcTextDecoder) _decompressGzip(bArr []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(bArr)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var output bytes.Buffer
	_, err = io.Copy(&output, reader)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// decode method
func (decoder *OmcTextDecoder) decode(file *os.File) ([]byte, error) {
	bytes, err := decoder.fileToByteArray(file)
	if err != nil {
		return nil, err
	}

	decoded := decoder._decode(bytes)
	decompressed, err := decoder._decompressGzip(decoded)
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}

// fileToByteArray method
func (decoder *OmcTextDecoder) fileToByteArray(file *os.File) ([]byte, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func main() {
	inputFile := os.Args[1]

	decoder := NewOmcTextDecoder()
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	decoded, err := decoder.decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(decoded))
	fmt.Println()
	fmt.Println("Decrypted xml has save to ./decypted_cscfeature.xml")
	os.WriteFile("./decypted_cscfeature.xml", decoded, 0644)
	// wait
	fmt.Println("Press 'Enter' to exit...")
	fmt.Scanln()
}
