package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

import . "particle"

func main() {
	var navg, nabsavg int = 0, 0
	var absmin, absavg float64 = 1.0, 0.0
	var davg, dmin float64 = 0, 0

	// Parse command line options
	np := flag.Int("n", 500, "sets the number of particles")
	op := flag.String("o", "out", "chooses and output file name")
	flag.Parse()

	// Set based on terminal
	var n int = *np
	var savename string = *op

	// Open file
	savefile, _ := os.Create(savename)
	defer savefile.Close()

	savewriter := bufio.NewWriter(savefile)
	defer savewriter.Flush()

	particles := make([]Particle, n)

	Set_size(n)
	Init_particles(n, particles)

	simulation_time := Get_time()
	for step := 0; step < NSTEPS; step++ {
		navg = 0
		davg = 0.0
		dmin = 1.0

		// Compute forces
		for i := 0; i < n; i++ {
			Zero_acceleration(&particles[i])

			for j := 0; j < n; j++ {
				Apply_force(&particles[i], &particles[j], &dmin, &davg, &navg)
			}
		}

		// Move particles
		for i := 0; i < n; i++ {
			Move(&particles[i])
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
