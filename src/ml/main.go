package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"./mutant"
	"./stats"
)

type isMutantRequest struct {
	Dna []string `json:"dna"`
}

type statsResponse struct {
	Mutant   int     `json:"count_mutant_dna"`
	NoMutant int     `json:"count_human_dna"`
	Ratio    float64 `json:"ratio"`
}

func handleStats(w http.ResponseWriter, req *http.Request) {
	mutant, noMutant := stats.Stats()
	ratio := float64(mutant) / math.Max(float64(1), float64(noMutant))
	response := statsResponse{
		Mutant:   mutant,
		NoMutant: noMutant,
		Ratio:    ratio,
	}
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=5")
	w.Write(data)
}

func handleIsMutatant(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var input isMutantRequest
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Set("Cache-Control", "max-age=31104000") // 1 year of cache
		ret, err := mutant.IsMutant(input.Dna)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			stats.StoreResult(input.Dna, ret)
			if ret {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		}
	}
}

func main() {
	http.HandleFunc("/mutant", handleIsMutatant)
	http.HandleFunc("/stats", handleStats)
	fmt.Print("server started...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
