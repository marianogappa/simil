package main

import (
	"math"
	"math/rand"
)

type centroid struct {
	center []float64
	ps     [][]float64
	is     []int
}

func d(p1, p2 []float64) float64 {
	sum := float64(0)
	for e := 0; e < len(p1); e++ {
		sum += math.Pow(p1[e]-p2[e], 2)
	}
	return math.Sqrt(sum)
}

func (c *centroid) recenter() float64 {
	newC := make([]float64, len(c.center))
	for _, e := range c.ps {
		for r := 0; r < len(newC); r++ {
			newC[r] += e[r]
		}
	}
	for r := 0; r < len(newC); r++ {
		newC[r] /= float64(len(c.ps))
	}
	oldCenter := c.center
	c.center = newC
	return d(oldCenter, c.center)
}

func kmeans(data [][]float64, k uint64, deltaThreshold float64) (centroids []centroid) {
	for i := uint64(0); i < k; i++ {
		centroids = append(centroids, centroid{center: data[rand.Intn(len(data))]})
	}

	for {
		for i := range data {
			minDist := math.MaxFloat64
			z := 0
			for v, e := range centroids {
				dist := d(data[i], e.center)
				if dist < minDist {
					minDist, z = dist, v
				}
			}
			centroids[z].ps = append(centroids[z].ps, data[i])
			centroids[z].is = append(centroids[z].is, i)
		}
		maxDelta := -math.MaxFloat64
		for i := range centroids {
			delta := centroids[i].recenter()
			if delta > maxDelta {
				maxDelta = delta
			}
		}
		if deltaThreshold >= maxDelta {
			return
		}
		for i := range centroids {
			centroids[i].ps = nil
		}
	}
}
