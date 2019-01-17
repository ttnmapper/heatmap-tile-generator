package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"log"
	"math"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawGatewayCountTile(x int, y int, z int, entries []types.MysqlAggGridcell) {

	tileNW := gosm.NewTileWithXY(x, y, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	//if pixelsPer19Tile < 1 {
	//	log.Printf("Level19 tileNW is less than one pixel")
	//}

	// fill in matrix here
	dc := gg.NewContext(256, 256)

	points := make(map[int]map[int]int)

	for _, entry := range entries {
		if _, ok := points[entry.X]; !ok {
			points[entry.X] = make(map[int]int)
		}
		if _, ok := points[entry.X][entry.Y]; !ok {
			points[entry.X][entry.Y] = 0
		}
		points[entry.X][entry.Y] += 1
	}

	for x, _ := range points {
		for y, _ := range points[x] {
			minRadius := 3.0
			nominalRadius := math.Max(minRadius, pixelsPer19Tile)

			pixelX := float64(x-tileNW19.X) * pixelsPer19Tile
			pixelY := float64(y-tileNW19.Y) * pixelsPer19Tile

			dc.DrawRectangle(pixelX, pixelY, nominalRadius, nominalRadius)

			switch points[x][y] {
			case 0:
				dc.SetRGB(0, 0, 0)
			case 1:
				dc.SetRGB255(255, 255, 212)
			case 2:
				dc.SetRGB255(254, 217, 142)
			case 3:
				dc.SetRGB255(254, 153, 41)
			default:
				dc.SetRGB255(204, 76, 2)
			}
			dc.Fill()
		}
	}

	// Write to file
	tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirGatewayCount, z, x)
	CreateDirIfNotExist(tilePath)
	tilePath = fmt.Sprintf("%s/%d.png", tilePath, y)
	err := dc.SavePNG(tilePath)
	if err != nil {
		log.Print(err.Error())
	}
}
