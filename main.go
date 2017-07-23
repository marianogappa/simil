package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type freq struct {
	i, c int
}

func generateFreqs(in []string) map[string]freq {
	var (
		fs = make(map[string]freq)
		j  = -1
	)
	for _, i := range in {
		for _, w := range strings.Split(i, " ") {
			if _, ok := fs[w]; !ok {
				j++
				fs[w] = freq{j, 1}
			} else {
				t := fs[w]
				t.c++
				fs[w] = t
			}
		}
	}
	return fs
}

// features are:
// 1. [] frequency one hot vector
// 2. sentence length
// 3. index of first word in wordmap
// 4. index of last word in wordmap
// 5. [] index of each word in sentence
func generateOneHots(in []string, fs map[string]freq) [][]float64 {
	var ohs = make([][]float64, len(in))
	for i := range ohs {
		ohs[i] = make([]float64, len(fs)+3)
	}

	var maxWordCount = 0
	for i := range in {
		var split = strings.Fields(in[i])
		for j, w := range split {
			var f = fs[w]
			ohs[i][f.i] = float64(f.c) // 1.
			if j == 0 {
				ohs[i][len(fs)+1] = float64(fs[w].i) // 3.
			}
			if j == len(split)-1 {
				ohs[i][len(fs)+2] = float64(fs[w].i) // 4.
			}
		}
		if len(split) > maxWordCount {
			maxWordCount = len(split)
		}
		ohs[i][len(fs)] = float64(len(in[i])) // 2.
	}

	for i := range in {
		var (
			split = strings.Split(in[i], " ")
			pos   = 0
		)
		for j := 0; j < maxWordCount; j++ {
			if len(split)-1 < j {
				ohs[i] = append(ohs[i], float64(0))
				continue
			}
			ohs[i] = append(ohs[i], float64(pos))
			pos += len(split[j]) + 1
		}
	}

	return ohs
}

func readInput(r io.Reader) []string {
	var (
		ls  = make([]string, 0, 500)
		rd  = bufio.NewReader(r)
		l   string
		err error
	)
	for {
		l, err = rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		ls = append(ls, l)
	}
	return ls
}

func out(short bool, i int, in string) {
	if short {
		fmt.Printf("%v\t%v\n", i, in)
	} else {
		fmt.Printf("cluster%v\t%v\n", i, in)
	}
}

func run(r io.Reader, k uint64, short, random bool, out func(bool, int, string)) {
	var (
		in  = readInput(r)
		fs  = generateFreqs(in)
		ohs = generateOneHots(in, fs)
		cs  = kmeans(ohs, k, 0.001)
	)
	for i := range cs {
		maxSumFs, maxSumFsI := 0, 0 // sentence in cluster with maximum aggregate frequency
		for _, p := range cs[i].ps {
			var inI int
		ohs:
			for j, oh := range ohs { // find the sentence (in[j]) from the feature vector (cs[i].ps)
				for l, c := range oh {
					if c != p[l] {
						continue ohs
					}
				}
				if !short {
					out(short, i, in[j])
				}
				inI = j
				break
			}

			if short {
				sumFs := 0 // calculate sentence aggregate frequency
				for _, f := range p[0:len(fs)] {
					sumFs += int(f)
				}
				if sumFs > maxSumFs { // find sentence with maximum aggregate frequency
					maxSumFs, maxSumFsI = sumFs, inI
				}
			}
		}
		if short && len(cs[i].ps) > 0 {
			out(short, len(cs[i].ps), in[maxSumFsI])
		}
	}
}

func main() {
	var (
		k      = flag.Uint64("k", 5, "how many clusters")
		short  = flag.Bool("short", false, "show one representative row per cluster, with cardinality")
		random = flag.Bool("random", false, "sets the random seed to current time (otherwise deterministic)")
	)

	flag.Parse()

	if *random {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	run(os.Stdin, *k, *short, *random, out)
}
