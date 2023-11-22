package main

import (
	"fmt"
	_ "fmt"
	_ "log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func start() {
	dat, err := os.ReadFile("input.txt")
	check(err)
	dataArray := strings.Fields(string(dat))

	dataSets, err := strconv.Atoi(dataArray[0])

	check(err)

	for i := 0; i < dataSets; i++ {

		numPoints, err := strconv.Atoi(dataArray[1])
		check(err)

		//remove this when creating the array of the representation
		degree, err := strconv.Atoi(dataArray[2])
		check(err)

		points := make([][]int, numPoints)
		for j := range points {
			points[j] = make([]int, 2)
		}

		for j := 0; j < numPoints; j++ {

			x, err := strconv.Atoi(dataArray[3+2*j])
			check(err)
			y, err := strconv.Atoi(dataArray[4+2*j])
			check(err)
			points[j][0] = x
			points[j][1] = y
		}

		fmt.Println(points)
		initialize(points, degree)

	}
}

func initialize(points [][]int, degree int) {
	rep := randFloats(-10, 10, degree+1)
	fmt.Println(rep)

}
func main() {
	start()
}
