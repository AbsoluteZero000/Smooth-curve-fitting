package main

import (
	"fmt"
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
		population := initialize(points, degree)
		parent1 := tournamentSelection(population)
		parent2 := tournamentSelection(population)

		// Do something with the selected parents...
	}
}

func initialize(points [][]int, degree int) [][]float64 {
	// Initialize the population and return it
	rep := randFloats(-10, 10, degree+1)
	fmt.Println(rep)

	// Return a population (you might want to have a population type)
	return [][]float64{rep}
}

func tournamentSelection(population [][]float64) []float64 {
	
	tournamentSize := 2					//tournment size b atnen
	selected := make([][]float64, tournamentSize)

	//  select random individuals ll mosb2a (tournament)
	for i := 0; i < tournamentSize; i++ {
		index := rand.Intn(len(population))
		selected[i] = population[index]
	}

	// a5tar a7sn wa7d mn l tournament
	bestIndividual := bestIndividual(selected)

	return bestIndividual
}
func bestIndividual(individuals [][]float64) []float64 {
	// nrg3 a7sn wa7d
	bestIndex := 0
	bestFitness := fitnessFunction(individuals[0])

	for i := 1; i < len(individuals); i++ {
		fitness := fitnessFunction(individuals[i])
		if fitness > bestFitness {
			bestFitness = fitness
			bestIndex = i
		}
	}
	return individuals[bestIndex]
}									 /////////////////////////////////////////////////////////////////////////////////
/*                                   /////////////////////////// 7d y3mlha ya shbab bmot ///////////////////////////
func fitnessFunction {				/////////////////////////////////////////////////////////////////////////////////
									/////////////////////////////////////////////////////////////////////////////////
}*/

func main() {
	start()
}
