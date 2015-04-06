// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Convenience functions in pure Go.
//
// Not all convenience functions are here: those that need
// to call C code have to be in core.go instead due to cgo
// limitations, while those that are methods have to be in
// core.go due to language limitations. They should all be
// here of course, at least conceptually.

package openal

import (
	"encoding/binary"
	"math"
	"os"
	"time"
	"fmt"
	"bytes"
)

var tempSlice = make([]float32, 6)

type WavFormat struct {
	FormatTag     int16
	Channels      int16
	Samples       int32
	AvgBytes      int32
	BlockAlign    int16
	BitsPerSample int16
}

type WavFormat2 struct {
	WavFormat
	SizeOfExtension int16
}

type WavFormat3 struct {
	WavFormat2
	ValidBitsPerSample int16
	ChannelMask        int32
	SubFormat          [16]byte
}

func ReadWavFile(path string) (*WavFormat, []byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var buff [4]byte
	f.Read(buff[:4])

	if string(buff[:4]) != "RIFF" {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	var size int32
	binary.Read(f, binary.LittleEndian, &size)

	f.Read(buff[:4])

	if string(buff[:4]) != "WAVE" {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	f.Read(buff[:4])

	if string(buff[:4]) != "fmt " {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	binary.Read(f, binary.LittleEndian, &size)

	var format WavFormat

	switch size {
	case 16:
		binary.Read(f, binary.LittleEndian, &format)
	case 18:
		var f2 WavFormat2
		binary.Read(f, binary.LittleEndian, &f2)
		format = f2.WavFormat
	case 40:
		var f3 WavFormat3
		binary.Read(f, binary.LittleEndian, &f3)
		format = f3.WavFormat
	}

	fmt.Println(format)

	f.Read(buff[:4])

	if string(buff[:4]) != "data" {
		return nil, nil, fmt.Errorf("Not supported WAV file.\n")
	}

	binary.Read(f, binary.LittleEndian, &size)

	wavData := make([]byte, size)
	n, e := f.Read(wavData)
	if e != nil {
		return nil, nil, fmt.Errorf("Cannot read WAV data.\n")
	}
	if int32(n) != size {
		return nil, nil, fmt.Errorf("WAV data size doesnt match.\n")
	}

	return &format, wavData, nil
}

func ReadWavFromMemory(memory []byte) (*WavFormat, []byte, error) {
	f := bytes.NewReader(memory)
	var buff [4]byte
	f.Read(buff[:4])

	if string(buff[:4]) != "RIFF" {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	var size int32
	binary.Read(f, binary.LittleEndian, &size)

	f.Read(buff[:4])

	if string(buff[:4]) != "WAVE" {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	f.Read(buff[:4])

	if string(buff[:4]) != "fmt " {
		return nil, nil, fmt.Errorf("Not a WAV file.\n")
	}

	binary.Read(f, binary.LittleEndian, &size)

	var format WavFormat

	switch size {
	case 16:
		binary.Read(f, binary.LittleEndian, &format)
	case 18:
		var f2 WavFormat2
		binary.Read(f, binary.LittleEndian, &f2)
		format = f2.WavFormat
	case 40:
		var f3 WavFormat3
		binary.Read(f, binary.LittleEndian, &f3)
		format = f3.WavFormat
	}

	f.Read(buff[:4])

	if string(buff[:4]) != "data" {
		return nil, nil, fmt.Errorf("Not supported WAV file.\n")
	}

	binary.Read(f, binary.LittleEndian, &size)

	wavData := make([]byte, size)
	n, e := f.Read(wavData)
	if e != nil {
		return nil, nil, fmt.Errorf("Cannot read WAV data.\n")
	}
	if int32(n) != size {
		return nil, nil, fmt.Errorf("WAV data size doesnt match.\n")
	}

	return &format, wavData, nil
}

func Period(freq int, samples int) float64 {
	return float64(freq) * 2 * math.Pi * (1 / float64(samples))
}

func TimeToData(t time.Duration, samples int, channels int) int {
	return int((float64(samples)/(1/t.Seconds()))+0.5) * channels
}

func PlaySoundFromMemory(memory []byte) {
	device := OpenDevice("")
	context := device.CreateContext()
	context.Activate()

	source := NewSource()
	buffer := NewBuffer()

	format, data, err := ReadWavFromMemory(memory)
	if err != nil {
		panic(err)
	}

	switch format.Channels {
	case 1:
		buffer.SetData(FormatMono16, data[:len(data)], int32(format.Samples))
	case 2:
		buffer.SetData(FormatStereo16, data[:len(data)], int32(format.Samples))
	}

	source.SetBuffer(buffer)
	source.Play()
	for source.State() == Playing {
		// wait
	}

	context.Destroy()
}

// Convenience function, see GetInteger().
func GetDistanceModel() int32 {
	return getInteger(alDistanceModel)
}

// Convenience function, see GetFloat().
func GetDopplerFactor() float32 {
	return getFloat(alDopplerFactor)
}

// Convenience function, see GetFloat().
func GetDopplerVelocity() float32 {
	return getFloat(alDopplerVelocity)
}

// Convenience function, see GetFloat().
func GetSpeedOfSound() float32 {
	return getFloat(alSpeedOfSound)
}

// Convenience function, see GetString().
func GetVendor() string {
	return GetString(alVendor)
}

// Convenience function, see GetString().
func GetVersion() string {
	return GetString(alVersion)
}

// Convenience function, see GetString().
func GetRenderer() string {
	return GetString(alRenderer)
}

// Convenience function, see GetString().
func GetExtensions() string {
	return GetString(alExtensions)
}
