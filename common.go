package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Constants
var density float64 = 0.0005
var mass float64 = 0.01
var Cutoff float64 = 0.01
var min_r float64 = (Cutoff / 100.0)
var dt float64 = 0.0005

var NSTEPS int = 1000
var SAVEFREQ int = 10

// Global variable
var Size float64

var first = true

type Particle struct {
	X, Y, vx, vy, ax, ay float64
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Find_index(bin []*Particle, p *Particle) int {
	for i := range bin {
		if bin[i] == p {
			return i
		}
	}
	return -1
}

func Delete_at_index(bin []*Particle, ind int) []*Particle {
	bin[ind] = bin[len(bin)-1]
	bin[len(bin)-1] = nil
	return bin[:len(bin)-1]
}

func Zero_acceleration(particle *Particle) {
	particle.ax = 0
	particle.ay = 0
}

func Get_time() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Second)
}

func Set_size(n int) {
	Size = math.Sqrt(float64(n) * density)
}

func Init_particles(n int, particles []Particle) {
	rand.Seed(time.Now().UTC().UnixNano())

	var sx int = int(math.Ceil(math.Sqrt(float64(n))))
	var sy int = (n + sx - 1) / sx

	shuffle := make([]int, n)
	for i := 0; i < n; i++ {
		shuffle[i] = i
	}

	for i := 0; i < n; i++ {

		// make sure particles are not spatially sorted
		var j int = rand.Intn(int(^uint32(1<<31))) % (n - i)
		var k int = shuffle[j]
		shuffle[j] = shuffle[n-i-1]

		// distribute particles evenly to ensure proper spacing
		particles[i].X = Size * float64(1+(k%sx)) / float64(1+sx)
		particles[i].Y = Size * float64(1+(k/sx)) / float64(1+sy)

		// assign random velocities within a bound
		particles[i].vx = rand.Float64()*2 - 1
		particles[i].vy = rand.Float64()*2 - 1
	}
}

func Apply_force(particle, neighbor *Particle, dmin, davg *float64, navg *int) {
	dx := neighbor.X - particle.X
	dy := neighbor.Y - particle.Y
	r2 := dx*dx + dy*dy

	if r2 > Cutoff*Cutoff {
		return
	}

	if r2 != 0 {
		if r2/(Cutoff*Cutoff) < (*dmin)*(*dmin) {
			*dmin = math.Sqrt(r2) / Cutoff
			(*davg) += math.Sqrt(r2) / Cutoff
			(*navg)++
		}
	}

	r2 = math.Max(r2, min_r*min_r)
	var r float64 = math.Sqrt(r2)

	var coeff float64 = (1.0 - Cutoff/r) / r2 / mass
	particle.ax += coeff * dx
	particle.ay += coeff * dy
}

func Move(p *Particle) {
	p.vx += p.ax * dt
	p.vy += p.ay * dt
	p.X += p.vx * dt
	p.Y += p.vy * dt

	for p.X < 0 || p.X > Size {
		if p.X < 0 {
			p.X = -p.X
		} else {
			p.X = 2*Size - p.X
		}
		p.vx = -p.vx
	}

	for p.Y < 0 || p.Y > Size {
		if p.Y < 0 {
			p.Y = -p.Y
		} else {
			p.Y = 2*Size - p.Y
		}
		p.vy = -p.vy
	}
}

func Save(w *bufio.Writer, n int, p []Particle) {
	if first {
		fmt.Fprintf(w, "%d %g\n", n, Size)
		first = false
	}

	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "%g %g\n", p[i].X, p[i].Y)
	}
}
