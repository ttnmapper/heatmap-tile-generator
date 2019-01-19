package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/j4/gosm"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"ttnmapper-heatmap-tile-generator/types"
)

func drawGlobalTile(x int, y int, z int, entries []types.MysqlAggGridcell) {

	// Our origin is one tile left and up, because we process 3x3 tiles
	tileNW := gosm.NewTileWithXY(x-1, y-1, z)
	tileNW19 := gosm.NewTileWithLatLong(tileNW.Lat, tileNW.Long, 19)

	pixelsPer19Tile := 256 / (math.Pow(2, float64(19-z)))
	minRadius := 3.0
	nominalRadius := math.Max(minRadius, pixelsPer19Tile)

	// fill in matrix here
	dc := gg.NewContext(768, 768)

	points := []types.Point{}

	for _, entry := range entries {
		newPoint := types.Point{}

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

	for _, entry := range points {

		pixelX := float64(entry.X-tileNW19.X) * pixelsPer19Tile
		pixelY := float64(entry.Y-tileNW19.Y) * pixelsPer19Tile

		// Move to centre of circle
		pixelX += (nominalRadius / 2)
		pixelY += (nominalRadius / 2)

		/*
			if($v['rssi'] < -120) $blue);
			elseif($v['rssi'] < -115) $cyan);
			elseif($v['rssi'] < -110) $green);
			elseif($v['rssi'] < -105) $yellow);
			elseif($v['rssi'] < -100) $orange);
			elseif($v['rssi'] < 0) $red);

			$blue = imagecolorallocatealpha($image, 0, 0, 255, 0);
			$cyan = imagecolorallocatealpha($image, 0, 255, 255, 0);
			$green = imagecolorallocatealpha($image, 0, 255, 0, 0);
			$yellow = imagecolorallocatealpha($image, 255, 255, 0, 0);
			$orange = imagecolorallocatealpha($image, 255, 127, 0, 0);
			$red = imagecolorallocatealpha($image, 255, 0, 0, 255);

			0 >-95  red 255 0 0
			1 >-100 red 255 0 0
			2 >-105 orange 255 127 0
			3 >-110 yellow 255 255 0
			4 >-115 green 0 255 0
			5 >-120 cyan 0 255 255
			6 >-125 blue 0 0 255
			7 >-130 blue
			8 >-135 blue
			9 >-140 blue
			10 >-145 blue
			11 <-145 blue
			12 - none black
		*/
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

	//// Write to file
	//tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirGlobalHeatmap, z, x)
	//CreateDirIfNotExist(tilePath)
	//tilePath = fmt.Sprintf("%s/%d-full.png", tilePath, y)
	//err := dc.SavePNG(tilePath)
	//if err != nil {
	//	log.Print(err.Error())
	//}

	srcImage := dc.Image()

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Crop out tile
			tile := srcImage.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(i*256, j*256, (i+1)*256, (j+1)*256))

			// Write to file
			tilePath := fmt.Sprintf("%s/%d/%d", myConfiguration.DirGlobalHeatmap, z, x-1+i)
			CreateDirIfNotExist(tilePath)
			tilePath = fmt.Sprintf("%s/%d.png", tilePath, y-1+j)

			newImage, _ := os.Create(tilePath)
			err := png.Encode(newImage, tile)
			if err != nil {
				log.Print(err.Error())
			}
			_ = newImage.Close()
		}
	}

}
