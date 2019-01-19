package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

var queueForToWrite = make(chan FileToWrite, 1000)

type FileToWrite struct {
	tile     image.Image
	dirName  string
	fileName string
}

func listenForFilesToWrite() {

	for {
		message := <-queueForToWrite
		writePNG(message)
	}
}

func writePNG(fileToWrite FileToWrite) {

	// Write to file
	CreateDirIfNotExist(fileToWrite.dirName)

	newImage, _ := os.Create(fileToWrite.dirName + "/" + fileToWrite.fileName)

	// TODO: PNG encoder is very slow
	//enc := &png.Encoder{
	//	CompressionLevel: png.BestSpeed,
	//}
	//err := enc.Encode(newImage, tile)
	err := png.Encode(newImage, fileToWrite.tile)
	if err != nil {
		log.Print(err.Error())
	}

	_ = newImage.Close()
}
