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

func fitnessFunction(coefficients []float64, points [][]float64) float64 {
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
func initialize(points [][]float64, degree int, popSize int) [][]float64 {

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
func tournamentSelection(population [][]float64, selectionSize int, points [][]float64) [][]float64 {

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

func bestIndividual(individuals [][]float64, points [][]float64) []float64 {
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

// /////////////////////////////// N-Point Cross Over //////////////////////////////////
func crossOver(selected [][]float64) [][]float64 {
	Pc := 0.75
	n := 2

	for i := 0; i < len(selected); i += 2 {
		rn := rand.Float64()
		if rn < Pc {

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

/////////////////////////////// Non-Uniform Mutation //////////////////////////////

func mutation(selected [][]float64, lowerBound int, upperBound int, generation int, maxGeneration int) [][]float64 {
	Pm := 0.05
	beta := 5.0

	for i := 0; i < len(selected); i++ {
		for j := 0; j < len(selected[i]); j++ {
			rn := rand.Float64()
			if rn < Pm {
				r := rand.Float64()
				delta := 0.0

				if r < 0.5 {
					delta = selected[i][j] - float64(lowerBound)
				} else {
					delta = float64(upperBound) - selected[i][j]
				}

				r = rand.Float64()
				delta *= (1 - math.Pow(r, math.Pow(1-float64(generation)/float64(maxGeneration), beta)))

				selected[i][j] = selected[i][j] + delta
			}
		}
	}

	return selected
}

/////////////////////// Elitist Replacement ////////////////////////

type TwoSlices struct {
	chromosomes [][]float64
	fitnesses   []float64
}

type SortByOther TwoSlices

func (sbo SortByOther) Len() int {
	return len(sbo.chromosomes)
}

func (sbo SortByOther) Swap(i, j int) {
	sbo.chromosomes[i], sbo.chromosomes[j] = sbo.chromosomes[j], sbo.chromosomes[i]
	sbo.fitnesses[i], sbo.fitnesses[j] = sbo.fitnesses[j], sbo.fitnesses[i]
}

func (sbo SortByOther) Less(i, j int) bool {
	return sbo.fitnesses[i] < sbo.fitnesses[j]
}

func replacement(generation [][]float64, copySize int, points [][]float64, offSpring [][]float64) [][]float64 {
	selected := make([][]float64, len(generation))
	fitnesses := make([]float64, len(generation))
	for i := 0; i < len(generation); i++ {
		fitnesses[i] = fitnessFunction(generation[i], points)
	}
	population := TwoSlices{chromosomes: generation, fitnesses: fitnesses}
	sort.Sort(SortByOther(population))
	j := 0
	for i := len(generation) - 1; i >= (len(generation) - copySize); i-- {
		selected[j] = population.chromosomes[i]
		j++
	}
	i := 0
	for ; j < len(generation); j++ {
		selected[j] = offSpring[i]
		i++
	}
	return selected
}

// //////////////////////// Main Function ////////////////////////////
func start() {
	maxGeneration := 1000
	popSize := 200
	selectionSize := 0.8 * float64(popSize)
	copiedParentsSize := popSize - int(selectionSize)
	lowerBound := -10
	upperBound := 10

	dat, err := os.ReadFile("input2.txt")
	check(err)
	dataArray := strings.Fields(string(dat))

	dataSets, err := strconv.Atoi(dataArray[0])
	check(err)
	degree := 0
	cursor := 0
	for i := 0; i < dataSets; i++ {
		cursor++
		numPoints, err := strconv.Atoi(dataArray[cursor])
		check(err)
		cursor++
		degree, err = strconv.Atoi(dataArray[cursor])
		check(err)

		points := make([][]float64, numPoints)
		for j := range points {
			points[j] = make([]float64, 2)
		}

		for j := 0; j < numPoints; j++ {
			x, err := strconv.ParseFloat(dataArray[3+2*j], 64)
			check(err)
			y, err := strconv.ParseFloat(dataArray[4+2*j], 64)
			check(err)
			points[j][0] = x
			points[j][1] = y
			cursor += 2
		}

		population := initialize(points, degree, popSize)

		for j := 0; j < maxGeneration; j++ {
			selectionPool := tournamentSelection(population, int(selectionSize), points)
			crossedOverPool := crossOver(selectionPool)
			mutatedPool := mutation(crossedOverPool, lowerBound, upperBound, i, maxGeneration)
			population = replacement(population, copiedParentsSize, points, mutatedPool)
		}
		fmt.Println("DATA SET", i+1, "\nBest Individual: ")
		for j := degree; j > 0; j-- {
			fmt.Print(population[0][j], " x^", j, " ", "+ ", " ")
		}
		fmt.Println(population[0][0])
	}
}

// /////////////////////////// Start ////////////////////////////////
func main() {
	rand.Seed(43)
	start()
}
