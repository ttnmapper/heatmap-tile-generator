package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"log"
	"math"
	"sort"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawPerGatewayTiles(x int, y int, z int, entries []types.MysqlAggGridcell) {

	tileNW := gosm.NewTileWithXY(x, y, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	//if pixelsPer19Tile < 1 {
	//	log.Printf("Level19 tileNW is less than one pixel")
	//}

	points := []types.Point{}

	for _, entry := range entries {
		newPoint := types.Point{}

		newPoint.GtwId = entry.GtwId

		newPoint.X = entry.X
		newPoint.Y = entry.Y

		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.BucketHigh)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket100)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket105)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket110)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket115)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket120)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket125)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket130)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket135)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket140)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.Bucket145)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.BucketLow)
		newPoint.BucketsValues = append(newPoint.BucketsValues, entry.BucketNoSignal)

		newPoint.MaxBucketIndex = 12

		for i := 0; i < len(newPoint.BucketsValues); i++ {
			if newPoint.BucketsValues[i] > newPoint.BucketsValues[newPoint.MaxBucketIndex] {
				newPoint.MaxBucketIndex = int8(i)
			}
		}

		points = append(points, newPoint)

	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].MaxBucketIndex > points[j].MaxBucketIndex
	})

	// fill in matrix here
	//dc := gg.NewContext(256, 256)
	images := make(map[string]*gg.Context)

	for _, entry := range points {
		minRadius := 16.0
		nominalRadius := math.Max(minRadius, pixelsPer19Tile)

		pixelX := float64(entry.X-tileNW19.X) * pixelsPer19Tile
		pixelY := float64(entry.Y-tileNW19.Y) * pixelsPer19Tile

		// Move to centre of circle
		pixelX += (nominalRadius / 2)
		pixelY += (nominalRadius / 2)

		if _, ok := images[entry.GtwId]; !ok {
			dc := gg.NewContext(256, 256)
			images[entry.GtwId] = dc
		}

		dc, _ := images[entry.GtwId]

		switch entry.MaxBucketIndex {
		case 0:
			dc.DrawCircle(pixelX, pixelY, nominalRadius)
			dc.SetRGB(1, 0, 0)
		case 1:
			dc.DrawCircle(pixelX, pixelY, nominalRadius+2)
			dc.SetRGB(1, 0.5, 0)
		case 2:
			dc.DrawCircle(pixelX, pixelY, nominalRadius+4)
			dc.SetRGB(1, 1, 0)
		case 3:
			dc.DrawCircle(pixelX, pixelY, nominalRadius+6)
			dc.SetRGB(0, 1, 0)
		case 4:
			dc.DrawCircle(pixelX, pixelY, nominalRadius+8)
			dc.SetRGB(0, 1, 1)
		case 5, 6, 7, 8, 9, 10, 11:
			dc.DrawCircle(pixelX, pixelY, nominalRadius+10)
			dc.SetRGB(0, 0, 1)
		case 12:
			dc.DrawCircle(pixelX, pixelY, nominalRadius)
			dc.SetRGB(0, 0, 0)
		}
		dc.Fill()
	}

	for gtwId, dc := range images {
		// Write to file
		tilePath := fmt.Sprintf("%s/%s/%d/%d", myConfiguration.DirGatewayHeatmap, gtwId, z, x)
		CreateDirIfNotExist(tilePath)
		tilePath = fmt.Sprintf("%s/%d.png", tilePath, y)
		err := dc.SavePNG(tilePath)
		if err != nil {
			log.Print(err.Error())
		}
	}
}
