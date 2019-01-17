package types

type Point struct {
	GtwId          string
	X              int
	Y              int
	BucketsValues  []int64
	MaxBucketIndex int8
}

type GatewayCountPoint struct {
	X     int
	Y     int
	Count int
}
