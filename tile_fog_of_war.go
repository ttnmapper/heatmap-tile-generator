package main

import (
	"fmt"
	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawFogOfWarTile(x int, y int, z int, entries []types.MysqlAggGridcell) {

	tileNW := gosm.NewTileWithXY(x, y, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	minRadius := 10.0
	nominalRadius := math.Max(minRadius, pixelsPer19Tile)

	// fill in matrix here
	dc := gg.NewContext(256, 256)

	for _, entry := range entries {

		pixelX := float64(entry.X-tileNW19.X) * pixelsPer19Tile
		pixelY := float64(entry.Y-tileNW19.Y) * pixelsPer19Tile

		// Move to centre of circle
		pixelX += (nominalRadius / 2)
		pixelY += (nominalRadius / 2)

		dc.DrawCircle(pixelX, pixelY, nominalRadius)
	}

	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, 256, 256)
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	srcImage := dc.Image()
	dstImage := image.NewRGBA(srcImage.Bounds())
	// Blur Function
	errBlur := graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: nominalRadius / 2.0})
	if errBlur != nil {
		log.Print(errBlur.Error())
	}

	// Write to file
	tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirFogOfWar, z, x)
	CreateDirIfNotExist(tilePath)
	tilePath = fmt.Sprintf("%s/%d.png", tilePath, y)

	newImage, _ := os.Create(tilePath)
	defer newImage.Close()
	err := png.Encode(newImage, dstImage)
	if err != nil {
		log.Print(err.Error())
	}
}
