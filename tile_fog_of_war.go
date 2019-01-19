package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"image"
	"log"
	"math"
	"time"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawFogOfWarTile(x int, y int, z int, entries []types.MysqlAggGridcell) {

	tileNW := gosm.NewTileWithXY(x-1, y-1, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	minRadius := 10.0
	nominalRadius := math.Max(minRadius, pixelsPer19Tile)

	// fill in matrix here
	dc := gg.NewContext(768, 768)

	for _, entry := range entries {

		pixelX := float64(entry.X-tileNW19.X) * pixelsPer19Tile
		pixelY := float64(entry.Y-tileNW19.Y) * pixelsPer19Tile

		dc.DrawCircle(pixelX, pixelY, nominalRadius)
	}

	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, 768, 768)
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	srcImage := dc.Image()
	// Blur Function

	start := time.Now()

	// Faster if we make the size smaller (default is 6 x stdDev)
	//blurredImage := image.NewRGBA(srcImage.Bounds())
	//errBlur := graphics.Blur(blurredImage, srcImage, &graphics.BlurOptions{StdDev: nominalRadius / 2.0, Size: int(nominalRadius)})
	//if errBlur != nil {
	//	log.Print(errBlur.Error())
	//}

	blurredImage := imaging.Blur(srcImage, nominalRadius/2.0)

	elapsed := time.Since(start)
	log.Printf("    blur took %s", elapsed)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Crop out tile
			tile := blurredImage.SubImage(image.Rect(i*256, j*256, (i+1)*256, (j+1)*256))

			// Write to file
			tileDirName := fmt.Sprintf("%s/%d/%d", myConfiguration.DirFogOfWar, z, x-1+i)
			tileFileName := fmt.Sprintf("%d.png", y-1+j)
			queueForToWrite <- FileToWrite{tile: tile, dirName: tileDirName, fileName: tileFileName}
		}
	}
}
