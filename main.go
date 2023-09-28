package main

import (
	"fmt"
	"math"
)

type UserPoint struct {
	Coordinates Point
	Status      string
}

type Point struct {
	Latitude  float64
	Longitude float64
}

type DBSCAN struct {
	UserPoints []UserPoint
	Epsilon    float64
	Clusters   []Cluster
	MinPoints  int
}

type Cluster struct {
	Points      []UserPoint
	CountPoints int
	UpperLeft   Point
	LowerRight  Point
	Center      Point
}

func Distance(p1, p2 UserPoint) float64 {
	latDiff := p1.Coordinates.Latitude - p2.Coordinates.Latitude
	lonDiff := p1.Coordinates.Longitude - p2.Coordinates.Longitude
	return math.Sqrt(latDiff*latDiff + lonDiff*lonDiff)
}

func (dbscan *DBSCAN) GetNeighbors(originalPoint UserPoint, neighbors []UserPoint) []UserPoint {
	if neighbors == nil {
		neighbors = make([]UserPoint, 0)
	}
	for index, point := range dbscan.UserPoints {
		currentPoint := &dbscan.UserPoints[index]
		if currentPoint.Status == "" && Distance(*currentPoint, originalPoint) <= dbscan.Epsilon {
			neighbors = append(neighbors, point)
			currentPoint.Status = "Attended"
			neighbors = dbscan.GetNeighbors(*currentPoint, neighbors)
		}
	}
	return neighbors
}

func (dbscan *DBSCAN) GetClusters() {
	for _, point := range dbscan.UserPoints {
		if point.Status == "" {
			points := dbscan.GetNeighbors(point, nil)
			if len(points) < dbscan.MinPoints {
				point.Status = "Visited"
				continue
			}

			cluster := Cluster{
				Points:      points,
				CountPoints: len(points),
			}
			cluster.GetCoordinatesCluster()
			cluster.GetClusterCenter()
			dbscan.Clusters = append(dbscan.Clusters, cluster)
		}
	}
}
func (cluster *Cluster) GetCoordinatesCluster() {
	maxLatitude := -90.0
	minLatitude := 90.0
	maxLongitude := -180.0
	minLongitude := 180.0
	for _, point := range cluster.Points {
		if point.Coordinates.Latitude > maxLatitude {
			maxLatitude = point.Coordinates.Latitude
		}
		if point.Coordinates.Latitude < minLatitude {
			minLatitude = point.Coordinates.Latitude
		}
		if point.Coordinates.Longitude > maxLongitude {
			maxLongitude = point.Coordinates.Longitude
		}
		if point.Coordinates.Longitude < minLongitude {
			minLongitude = point.Coordinates.Longitude
		}
	}

	cluster.UpperLeft = Point{
		Latitude:  maxLatitude,
		Longitude: minLongitude,
	}
	cluster.LowerRight = Point{
		Latitude:  minLatitude,
		Longitude: maxLongitude,
	}
}

func (cluster *Cluster) GetClusterCenter() {
	cluster.Center = Point{
		Latitude:  (cluster.UpperLeft.Latitude + cluster.LowerRight.Latitude) / 2,
		Longitude: (cluster.LowerRight.Longitude + cluster.UpperLeft.Longitude) / 2,
	}

}

func main() {
	points := []UserPoint{
		{Coordinates: Point{Latitude: 40.7128, Longitude: -74.0060}},
		{Coordinates: Point{Latitude: 40.7128, Longitude: -74.0060}},
		{Coordinates: Point{Latitude: 40.9128, Longitude: -74.0060}},
		{Coordinates: Point{Latitude: 43.9128, Longitude: -74.0060}},
		{Coordinates: Point{Latitude: 44.9128, Longitude: -75.0060}},
		{Coordinates: Point{Latitude: 54.9129, Longitude: -75.0060}},
		{Coordinates: Point{Latitude: 64.9128, Longitude: -75.0060}},
	}
	dbscan := DBSCAN{Epsilon: 10.01, UserPoints: points, MinPoints: 1}
	dbscan.GetClusters()

	for _, cluster := range dbscan.Clusters {
		fmt.Printf("Cluster:\n")
		fmt.Printf("Count of points: %d\n", cluster.CountPoints)
		fmt.Printf("Upper left corner: (%f, %f)\n", cluster.UpperLeft.Latitude, cluster.UpperLeft.Longitude)
		fmt.Printf("Lower right corner: (%f, %f)\n", cluster.LowerRight.Latitude, cluster.LowerRight.Longitude)
		fmt.Printf("Center: (%f, %f)\n", cluster.Center.Latitude, cluster.Center.Longitude)
		fmt.Printf("Points:\n")
		for _, point := range cluster.Points {
			fmt.Printf("(%f, %f)\n", point.Coordinates.Latitude, point.Coordinates.Longitude)
		}
		fmt.Println()
	}

}
