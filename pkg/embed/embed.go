package embed

import (
	"math"
	"strings"
	"sync"
	"unicode"
)

type Vectorizer struct {
	mu       sync.RWMutex
	docCount int
	docFreq  map[string]int
}

func NewVectorizer() *Vectorizer {
	return &Vectorizer{
		docFreq: make(map[string]int),
	}
}

func tokenize(text string) []string {
	var words []string
	var b strings.Builder
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			if b.Len() > 0 {
				words = append(words, b.String())
				b.Reset()
			}
		}
	}
	if b.Len() > 0 {
		words = append(words, b.String())
	}
	return words
}

func (v *Vectorizer) Train(documents []string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.docCount += len(documents)
	for _, doc := range documents {
		tokens := tokenize(doc)
		seen := make(map[string]bool)
		for _, t := range tokens {
			if !seen[t] {
				seen[t] = true
				v.docFreq[t]++
			}
		}
	}
}

func (v *Vectorizer) Transform(text string) map[string]float64 {
	v.mu.RLock()
	defer v.mu.RUnlock()

	tokens := tokenize(text)
	if len(tokens) == 0 {
		return nil
	}

	tf := make(map[string]float64)
	for _, t := range tokens {
		tf[t] += 1.0
	}

	norm := 0.0
	vec := make(map[string]float64)
	totalDocs := float64(v.docCount)
	if totalDocs == 0 {
		totalDocs = 1.0
	}

	for term, count := range tf {
		df := float64(v.docFreq[term])
		if df == 0 {
			df = 1.0
		}
		idf := math.Log((totalDocs + 1.0) / (df + 1.0)) + 1.0
		weight := (count / float64(len(tokens))) * idf
		vec[term] = weight
		norm += weight * weight
	}

	norm = math.Sqrt(norm)
	if norm > 0 {
		for term := range vec {
			vec[term] /= norm
		}
	}

	return vec
}

func CosineSimilarity(a, b map[string]float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	dot := 0.0
	for k, va := range a {
		if vb, ok := b[k]; ok {
			dot += va * vb
		}
	}
	return dot
}
