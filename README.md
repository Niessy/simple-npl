<h1>simple-nlp</h1>

Simple NLP library, will include PCFG, Translation IBM 1 and 2 Models, GLM, etc...

Done through nlp class on coursera

To use in your GOPATH, should be (YOUR GO DIR)/src/github.com/Niessy.
	
    git clone https://github.com/Niessy/simple-nlp.git
    
Or

    go get github.com/Niessy/simple-nlp

<h2> Probabilistic Context Free Grammar </h2>

Example of usage: This works assuming the counts file already has rare counts, I'm going to add support for this soon.

    package main

    import (
        "fmt"
	    "github.com/Niessy/simple-nlp/pcfg"
    )

    func main() {
	    p := pcfg.NewPCFG(5)
	    p.GetCounts("counts.test")
	    err := p.ParseSentences("parse_dev.dat", "output.json")
	    if err != nil {
		    fmt.Println(err)
	    }
    }
    
<h2> Machine Translation </h2>
    
Example IBM1:

	package main
	
	import (
		"github.com/Niessy/simple-nlp/translation"
		"fmt"
	)
	
	func main() {
		// corpus.en is the native text, corpus.es acts as the foreign text
		// output.json contains associated probabilities for aligned native/foreign words
		// used as input for IBM2
		i := translation.NewIBM1("corpus.en", "corpus.es", "output.json")
		err := i.Initialize()
		if err != nil {
			fmt.Println(err)
		}
		err = i.EMAlgorithm(5) // Number of iterations is 5
		if err != nil {
			fmt.Println(err)
		}
	}
	
Example IBM2:

	package main

	import (
		"fmt"
		"github.com/Niessy/simple-nlp/translation"
	)
	
	func main() {
		i := translation.NewIBM2("output.json", "tparams.json", "qparams.json")
		err := i.Initialize("corpus.en", "corpus.es")
		if err != nil {
			fmt.Println(err)
		}
		err = i.EMAlgorithm(5)
		if err != nil {
			fmt.Println(err)
		}
	}

Example: Aligner, outputs the alignment of native/foreign words

	package main
	
	import (
		"fmt"
		"github.com/Niessy/simple-nlp/translation"
	)
	
	func main() {
		a := new(translation.Aligner)
		a.EnglishFile = "corpus.en"
		a.ForeignFile = "corpus.es"
		a.AlignmentFile = "alignment.txt"
		err = a.GetParams("tparams.json", "qparams.json")
		if err != nil {
			fmt.Println(err)
		}
		err = a.BestAlignment()
		if err != nil {
			fmt.Println(err)
		}
	}

    

