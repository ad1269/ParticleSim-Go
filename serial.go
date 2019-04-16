package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
)

import . "particle"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var navg, nabsavg int = 0, 0
	var absmin, absavg float64 = 1.0, 0.0
	var davg, dmin float64 = 0, 0

	// Parse command line options
	np := flag.Int("n", 1000, "sets the number of particles")
	op := flag.String("o", "out", "chooses and output file name")
	flag.Parse()

	// Set based on terminal
	var n int = *np
	var savename string = *op

	// Open file
	savefile, err := os.Create(savename)
	check(err)
	defer savefile.Close()

	savewriter := bufio.NewWriter(savefile)
	defer savewriter.Flush()

	particles := make([]Particle, n)

	Set_size(n)
	Init_particles(n, particles)

	// Set up binning
	bins := int(math.Ceil(Size / Cutoff))
	pointers := make([][]*Particle, bins*bins)
	for i := range pointers {
		pointers[i] = make([]*Particle, 0)
	}

	// Populate bins
	for i := 0; i < n; i++ {
		var bin_i int = Min(int(particles[i].X/Cutoff), bins-1)
		var bin_j int = Min(int(particles[i].Y/Cutoff), bins-1)

		pointers[bin_i*bins+bin_j] = append(pointers[bin_i*bins+bin_j], &particles[i])
	}

	simulation_time := Get_time()
	for step := 0; step < NSTEPS; step++ {
		navg = 0
		davg = 0.0
		dmin = 1.0

		// Compute forces
		for i := 0; i < n; i++ {
			Zero_acceleration(&particles[i])

			var bin_i int = Min(int(particles[i].X/Cutoff), bins-1)
			var bin_j int = Min(int(particles[i].Y/Cutoff), bins-1)

			for other_bin_i := Max(0, bin_i-1); other_bin_i <= Min(bins-1, bin_i+1); other_bin_i++ {
				for other_bin_j := Max(0, bin_j-1); other_bin_j <= Min(bins-1, bin_j+1); other_bin_j++ {
					for k := 0; k < len(pointers[other_bin_i*bins+other_bin_j]); k++ {
						var p *Particle = pointers[other_bin_i*bins+other_bin_j][k]
						Apply_force(&particles[i], p, &dmin, &davg, &navg)
					}
				}
			}
		}

		// Move particles
		for i := 0; i < n; i++ {
			var old_bin_i int = Min(int(particles[i].X/Cutoff), bins-1)
			var old_bin_j int = Min(int(particles[i].Y/Cutoff), bins-1)

			Move(&particles[i])

			var new_bin_i int = Min(int(particles[i].X/Cutoff), bins-1)
			var new_bin_j int = Min(int(particles[i].Y/Cutoff), bins-1)

			// Has been moved into a new bin
			if old_bin_j != new_bin_j || old_bin_i != new_bin_i {
				var old_bin []*Particle = pointers[old_bin_i*bins+old_bin_j]
				pointers[old_bin_i*bins+old_bin_j] = Delete_at_index(old_bin, Find_index(old_bin, &particles[i]))

				var new_bin []*Particle = pointers[new_bin_i*bins+new_bin_j]
				pointers[new_bin_i*bins+new_bin_j] = append(new_bin, &particles[i])
			}
		}

		// Compute Statistical data
		if navg != 0 {
			absavg += davg / float64(navg)
			nabsavg++
		}

		if dmin < absmin {
			absmin = dmin
		}

		if step%SAVEFREQ == 0 {
			Save(savewriter, n, particles)
		}
	}

	simulation_time = Get_time() - simulation_time
	fmt.Printf("n = %d, simulation time = %f seconds", n, simulation_time)

	if nabsavg != 0 {
		absavg /= float64(nabsavg)
	}

	//
	//  -The minimum distance absmin between 2 particles during the run of the simulation
	//  -A Correct simulation will have particles stay at greater than 0.4 (of cutoff) with typical values between .7-.8
	//  -A simulation where particles don't interact correctly will be less than 0.4 (of cutoff) with typical values between .01-.05
	//
	//  -The average distance absavg is ~.95 when most particles are interacting correctly and ~.66 when no particles are interacting
	//
	fmt.Printf(", absmin = %f, absavg = %f\n", absmin, absavg)
	if absmin < 0.4 {
		fmt.Println("The minimum distance is below 0.4 meaning that some particle is not interacting")
	}

	if absavg < 0.8 {
		fmt.Println("The average distance is below 0.8 meaning that most particles are not interacting")
	}
}
