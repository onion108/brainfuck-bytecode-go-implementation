package main

import (
	"bufio"
	"os"
)

const VMCodeLoaderMaxSize = 256

type VMCodeLoader struct {
	startAddr  int
	buffer     [VMCodeLoaderMaxSize]byte
	file       *os.File
	loaded     bool
	actualSize int
}

func VMCodeLoaderInit(file *os.File) *VMCodeLoader {
	return &VMCodeLoader{
		0,
		[VMCodeLoaderMaxSize]byte{},
		file,
		false,
		0,
	}
}

func (l *VMCodeLoader) loadFrom(addr int) {
	// Prepare
	l.file.Seek(int64(addr), 0)
	// Loading bytes
	// Create the reader
	reader := bufio.NewReader(l.file)
	for i := 0; i < VMCodeLoaderMaxSize; i++ {
		// Read a byte, add it into buffer
		currentByte, err := reader.ReadByte()
		// Read failed, exit
		if err != nil {
			l.actualSize = i
			break
		}
		l.buffer[i] = currentByte
	}
	// Set the actual size to the max because we read all the 256 bytes.
	l.actualSize = VMCodeLoaderMaxSize
	// End loading
}

func (l *VMCodeLoader) At(addr int) (byte, bool) {
	if addr < 0 {
		return 0xFF, false
	}
	// If it's unloaded
	if !l.loaded {
		l.loadFrom(addr / VMCodeLoaderMaxSize * VMCodeLoaderMaxSize)
		if addr >= l.startAddr+l.actualSize {
			// If out of bounds, return
			return 0xFF, false
		}
		// Return the data
		return l.buffer[addr%VMCodeLoaderMaxSize], true
	}
	// If it's loaded, but not in current loaded chunk
	if addr < l.startAddr || addr >= l.startAddr+l.actualSize {
		l.loadFrom(addr / VMCodeLoaderMaxSize * VMCodeLoaderMaxSize)
		if addr >= l.startAddr+l.actualSize {
			// If still out of bound,
			return 0xFF, false
		}
		return l.buffer[addr%VMCodeLoaderMaxSize], true
	}
	// Otherwise, it's loaded and in the current loaded chunk
	return l.buffer[addr%VMCodeLoaderMaxSize], true
}
