package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

// /////////////////////////////// Utility Functions //////////////////////////////
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func calculatePolynomial(x float64, coefficients []float64) float64 {
	sum := 0.0
	for i := 0; i < len(coefficients); i++ {
		sum += coefficients[i] * math.Pow(x, float64(i))
	}
	return sum
}

func fitnessFunction(coefficients []float64, points [][]int) float64 {
	totalError := 0.0
	N := float64(len(points))

	for _, point := range points {
		x := float64(point[0])
		yActual := float64(point[1])
		yCalc := calculatePolynomial(x, coefficients)
		error := math.Pow(yCalc-yActual, 2)

		totalError += error
	}

	mse := totalError / N

	fitness := 1.0 / mse

	return fitness
}

// /////////////////////////////// Population initialization //////////////////////////////////
func initialize(points [][]int, degree int, popSize int) [][]float64 {

	population := make([][]float64, popSize)
	for i := 0; i < popSize; i++ {
		individual := make([]float64, degree+1)

		for j := 0; j <= degree; j++ {
			individual[j] = -10 + rand.Float64()*(20)
		}

		population[i] = individual
	}

	return population
}

// /////////////////////////////// Tournament Selection /////////////////////////////////
func tournamentSelection(population [][]float64, selectionSize int, points [][]int) [][]float64 {

	tournamentSize := 2
	selected := make([][]float64, selectionSize)
	for i := 0; i < selectionSize; i++ {
		selected[i] = make([]float64, len(population[0]))
	}

	for j := 0; j < selectionSize; j++ {
		arena := make([][]float64, tournamentSize)
		for i := 0; i < tournamentSize; i++ {
			arena[i] = make([]float64, len(population[0]))

			index := rand.Intn(len(population))
			arena[i] = population[index]

		}

		selected[j] = bestIndividual(arena, points)
	}
	return selected
}

func bestIndividual(individuals [][]float64, points [][]int) []float64 {
	bestIndex := 0
	bestFitness := fitnessFunction(individuals[0], points)

	for i := 1; i < len(individuals); i++ {
		fitness := fitnessFunction(individuals[i], points)
		if fitness > bestFitness {
			bestFitness = fitness
			bestIndex = i
		}
	}
	return individuals[bestIndex]
}

// /////////////////////////////// Cross Over //////////////////////////////////
func crossOver(selected [][]float64) [][]float64 {
	Pc := 0.75
	n := 2

	for i := 0; i < len(selected); i += 2 {

		if rand.Float64() < Pc {
			crossoverPoints := make([]int, n)
			for j := 0; j < n; j++ {
				crossoverPoints[j] = rand.Intn(len(selected[i]))
			}
			sort.Ints(crossoverPoints)

			child1 := make([]float64, len(selected[i]))
			child2 := make([]float64, len(selected[i]))

			flip := false
			for j := 0; j < len(selected[i]); j++ {
				if flip {
					child1[j] = selected[i][j]
					child2[j] = selected[i+1][j]
				} else {
					child1[j] = selected[i+1][j]
					child2[j] = selected[i][j]
				}

				for k := 0; k < n-1; k++ {
					if j == crossoverPoints[k] {
						flip = !flip
						break
					}
				}
			}

			selected[i] = child1
			selected[i+1] = child2
		}
	}

	return selected
}

func start() {
	popSize := 50
	selectionSize := 10

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

		population := initialize(points, degree, popSize)

		selectionPool := tournamentSelection(population, selectionSize, points)

		fmt.Print(selectionPool)


	}
}

func main() { start() }
