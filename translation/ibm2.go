package translation

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type IBM2 struct {
	Ibm1        IBM1
	TparamInput string
	QparamsFile string
	Qparams     map[Q]float64
}

type Q struct {
	i int // foreign word
	j int // english word
	m int // foreign sentence length
	l int // english sentence length
}

func (model *IBM2) Initialize() error {
	// Read tparams generated by IBM1 from file,
	// this will probably be taken out later.
	// jusy gonna make sure program works first
	err := model.readtparams()
	if err != nil {
		return err
	}

	ff, err := os.Open(model.Ibm1.ForeignFile)
	if err != nil {
		return err
	}

	ef, err := os.Open(model.Ibm1.EnglishFile)
	if err != nil {
		return err
	}

	defer ff.Close()
	defer ef.Close()

	fr := bufio.NewReader(ff)
	er := bufio.NewReader(ef)

	for {
		engLine, err := er.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		forLine, err := fr.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if engLine == "\n" || forLine == "\n" {
			continue
		}

		engWords := strings.Split(strings.TrimSpace(engLine), " ")
		forWords := strings.Split(strings.TrimSpace(forLine), " ")
		el := len(engWords)
		ml := len(forWords)
		for i := 1; i <= ml; i++ {
			for j := 0; j <= el; j++ {
				q := Q{i: i, j: j, m: ml, l: el}
				model.Qparams[q] = (1.0 / (float64(el) + 1))
			}
		}
	}
	return nil
}

func (model *IBM2) EMAlgorithm(n int) error {
	tt := time.Now()
	fmt.Printf("Beginning EM for IBM2\n")
	for s := 1; s <= n; s++ {
		fmt.Printf("Beginning iteration %d...\n", s)
		t1 := time.Now()
		// Reset counts
		counts := make(map[string]float64)
		qcounts := make(map[Q]float64)

		// being reading files
		ff, err := os.Open(model.Ibm1.ForeignFile)
		if err != nil {
			return err
		}

		ef, err := os.Open(model.Ibm1.EnglishFile)
		if err != nil {
			return err
		}

		er := bufio.NewReader(ef)
		fr := bufio.NewReader(ff)

		for {
			engLine, err := er.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			forLine, err := fr.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if engLine == "\n" || forLine == "\n" {
				continue
			}

			engWords := strings.Split("NULL "+strings.TrimSpace(engLine), " ")
			forWords := strings.Split(strings.TrimSpace(forLine), " ")

			i := 1
			el := len(engWords) - 1
			fl := len(forWords)
			for _, fw := range forWords {
				for j, ew := range engWords {
					d := model.delta2(i, j, el, fl, ew, fw, &engWords)
					counts[ew+" "+fw] = counts[ew+" "+fw] + d
					counts[ew] = counts[ew] + d

					q := Q{i: i, j: j, m: fl, l: el}
					q2 := Q{i: -1, j: -1, m: fl, l: el}

					qcounts[q] = qcounts[q] + d
					qcounts[q2] = qcounts[q2] + d
				}
				i++
			}

		}
		// Cleanup files for next iteration
		ff.Close()
		ef.Close()

		// Revising Tparams
		for e, fws := range model.Ibm1.Tparams {
			for f, _ := range fws {
				model.Ibm1.Tparams[e][f] = counts[e+" "+f] / counts[e]
			}
		}

		for key, _ := range model.Qparams {
			top := qcounts[key]
			model.Qparams[key] = top / qcounts[Q{i: -1, j: -1, m: key.m, l: key.l}]
		}

		fmt.Printf("Ending iteration %d...took %v\n", s, time.Since(t1))
	}
	fmt.Println("Finished EM Algorithm took ", time.Since(tt))
	fmt.Printf("\nWriting tparams to %s...\n", model.Ibm1.TparamsFile)
	err := model.Ibm1.tparamsToFile()
	if err != nil {
		return err
	}
	fmt.Printf("Finished writing tparams...\n")
	fmt.Printf("\nWriting qparams to %s...\n", model.QparamsFile)
	err = model.qparamsToFile()
	if err != nil {
		return err
	}
	fmt.Printf("Finished writing qparams...\n")
	return nil
}
