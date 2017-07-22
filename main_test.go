package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

type settings struct {
	K      uint64
	Short  bool
	Random bool
}

func TestSimil(t *testing.T) {
	dir := "./testdata"
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".settings") {
			name := f.Name()[:len(f.Name())-9]
			var settings settings

			r, err := os.Open(dir + "/" + name + ".in")
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			defer r.Close()

			byts, err := ioutil.ReadFile(dir + "/" + name + ".settings")
			if err != nil {
				t.Errorf("Couldn't read settings file. err=%v", err)
				t.FailNow()
			}
			err = json.Unmarshal(byts, &settings)
			if err != nil {
				t.Errorf("Couldn't unmarshal settings file. err=%v", err)
				t.FailNow()
			}

			expected := [][]string{}
			file, err := os.Open(dir + "/" + name + ".expected")
			if err != nil {
				t.Errorf("Couldn't read expected file. err=%v", err)
			}
			defer file.Close()

			var rd = bufio.NewReader(file)
			var line, curCl string
			for {
				line, err = rd.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Errorf("Couldn't read line from expected file. err=%v", err)
					t.FailNow()
				}
				line := line[:len(line)-1]
				fields := strings.Split(line, "\t")
				if len(fields) != 2 {
					t.Errorf("Line should have 2 fields; instead: %+v", fields)
					t.FailNow()
				}
				if curCl == "" {
					curCl = fields[0]
					expected = append(expected, []string{})
				}
				if curCl != fields[0] {
					curCl = fields[0]
					expected = append(expected, []string{})
				}
				expected[len(expected)-1] = append(expected[len(expected)-1], fields[1])
			}

			ai := -1
			actual := [][]string{}
			testOut := func(short bool, i int, in string) {
				if ai == -1 {
					ai = i
					actual = append(actual, []string{})
				}
				if ai != i {
					ai = i
					actual = append(actual, []string{})
				}
				actual[len(actual)-1] = append(actual[len(actual)-1], in)
			}

			run(bufio.NewReader(r), settings.K, false, settings.Random, testOut)

			compare(t, expected, actual, 0.8)
		}
	}
}

func compare(t *testing.T, es [][]string, as [][]string, correctness float64) {
	for _, e := range es {
		sort.Strings(e)
		ok := false
		for _, a := range as {
			sort.Strings(a)
			ei, ai, matches := 0, 0, 0
			for {
				if ei >= len(e) || ai >= len(a) {
					break
				}
				if e[ei] == a[ai] {
					matches++
					ei++
					ai++
					continue
				}
				if lt(e[ei], a[ai]) {
					ei++
				} else {
					ai++
				}
			}
			if float64(matches)/float64(len(e)) >= correctness {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("No match for cluster %+v", e)
			t.FailNow()
		}
	}
}

func lt(a string, b string) bool {
	return sort.StringsAreSorted([]string{a, b})
}
