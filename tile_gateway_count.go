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

	// Aggregate number of gateways per z19 cell
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

	// z19 cells to pixels
	blockSize := math.Max(pixelsPer19Tile, 16)

	matrix := make(map[int]map[int]int)
	for x, _ := range points {
		for y, _ := range points[x] {

			pointX := float64(x-tileNW19.X) * pixelsPer19Tile
			pointY := float64(y-tileNW19.Y) * pixelsPer19Tile

			if pointX < 0 || pointX > 255 || pointY < 0 || pointY > 255 {
				continue
			}

			blockX := int(pointX / blockSize)
			blockY := int(pointY / blockSize)

			if _, ok := matrix[blockX]; !ok {
				matrix[blockX] = make(map[int]int)
			}
			if _, ok := matrix[blockX][blockY]; !ok {
				matrix[blockX][blockY] = 0
			}

			if matrix[blockX][blockY] < points[x][y] {
				matrix[blockX][blockY] = points[x][y]
			}

		}
	}

	dc := gg.NewContext(256, 256)

	for x, _ := range matrix {
		for y, _ := range matrix[x] {

			dc.DrawRectangle(float64(x)*blockSize, float64(y)*blockSize, blockSize, blockSize)

			switch matrix[x][y] {
			case 0:
				dc.SetRGBA(0, 0, 0, 0)
			case 1:
				dc.SetRGB(0, 0, 1)
			case 2:
				dc.SetRGB(0, 0.5, 0)
			case 3:
				dc.SetRGB(1, 0.5, 0)
			default:
				dc.SetRGB(1, 0, 0)
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

	//srcImage := dc.Image()
	//
	//for i := 0; i<3; i++ {
	//	for j := 0; j<3; j++ {
	//		// Crop out tile
	//		tile := srcImage.(interface {
	//			SubImage(r image.Rectangle) image.Image
	//		}).SubImage(image.Rect(i*256, j*256, (i+1)*256, (j+1)*256))
	//
	//		// Write to file
	//		tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirGatewayCount, z, x-1+i)
	//		CreateDirIfNotExist(tilePath)
	//		tilePath = fmt.Sprintf("%s/%d.png", tilePath, y-1+j)
	//
	//		newImage, _ := os.Create(tilePath)
	//		err := png.Encode(newImage, tile)
	//		if err != nil {
	//			log.Print(err.Error())
	//		}
	//		_ = newImage.Close()
	//	}
	//}
}
