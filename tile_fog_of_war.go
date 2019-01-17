package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"log"
	"math"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawFogOfWarTile(x int, y int, z int, entries []types.MysqlAggGridcell) {

	tileNW := gosm.NewTileWithXY(x, y, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	//if pixelsPer19Tile < 1 {
	//	log.Printf("Level19 tileNW is less than one pixel")
	//}

	// fill in matrix here
	dc := gg.NewContext(256, 256)
	//dc.DrawRectangle(0, 0, 256, 256)
	//dc.SetRGBA(0, 0, 0, 1)
	//dc.Fill()

	for _, entry := range entries {
		minRadius := 3.0
		nominalRadius := math.Max(minRadius, pixelsPer19Tile)

		pixelX := float64(entry.X-tileNW19.X) * pixelsPer19Tile
		pixelY := float64(entry.Y-tileNW19.Y) * pixelsPer19Tile

		// Move to centre of circle
		pixelX += (nominalRadius / 2)
		pixelY += (nominalRadius / 2)

		dc.DrawCircle(pixelX, pixelY, nominalRadius)
		//dc.SetRGBA(1, 1, 1, 1)
		//dc.Fill()
	}

	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, 256, 256)
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	// Write to file
	tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirFogOfWar, z, x)
	CreateDirIfNotExist(tilePath)
	tilePath = fmt.Sprintf("%s/%d.png", tilePath, y)
	err := dc.SavePNG(tilePath)
	if err != nil {
		log.Print(err.Error())
	}
}
