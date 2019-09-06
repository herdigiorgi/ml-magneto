package mutant

import "testing"

func assertMutantError(t *testing.T, input []string) {
	_, err := IsMutant(input)
	if err == nil {
		t.Errorf("error expected for %v", input)
	}
}

func assertMutant(t *testing.T, input []string) {
	res, err := IsMutant(input)
	if err != nil {
		t.Errorf("%v with error `%v`", input, err)
	} else if res == false {
		t.Errorf("%v isn't a mutant", input)
	}
}

func assertNormal(t *testing.T, input []string) {
	res, err := IsMutant(input)
	if err != nil {
		t.Errorf("%v with error `%v`", input, err)
	} else if res {
		t.Errorf("%v should be normal", input)
	}
}

func TestIsMutantSmoke(t *testing.T) {
	IsMutant([]string{"ATGCGA", "ATGCGA"})
}

func TestIsMutantInvalidInput(t *testing.T) {
	assertMutantError(t, []string{"ATGCGA", "APGCGA"})
	assertMutantError(t, []string{"ATGCGA", "ATGCG"})
	assertMutantError(t, []string{"ATGCG", "ATGCGA"})
}

func TestIsMutantRow(t *testing.T) {
	assertMutant(t, []string{"AAAA"})
	assertMutant(t, []string{"TAAAA"})
	assertMutant(t, []string{
		"TGCG",
		"AAAA", // AAAA
	})
	assertMutant(t, []string{
		"AAGCCA",
		"ATGGGG", // --GGGG
		"AAGTGA",
	})
	assertMutant(t, []string{
		"AAGCCA",
		"ATGGGG", // --GGGG
		"AAGTGA",
	})
	assertMutant(t, []string{
		"AAGCCA",
		"ATGGCG",
		"TAAAAG", // -AAAA-
	})
}

func TestIsMutantColumn(t *testing.T) {
	assertNormal(t, []string{"A", "A", "A"})
	assertMutant(t, []string{"A", "A", "A", "A"})
	assertMutant(t, []string{
		"TAG", // A
		"GAC", // A
		"TAG", // A
		"GAC", // A
	})
	assertMutant(t, []string{
		"CTG", // ---
		"TAG", // -A-
		"GAC", // -A-
		"TAG", // -A-
		"GAC", // -A-
		"GTC", // ---
	})
}

func TestIsMutantColumnOblique(t *testing.T) {
	assertMutant(t, []string{
		"CTGTT", // -----
		"TGCGA", // -G---
		"GAGTA", // --G--
		"TAGGT", // ---G-
		"GACGG", // ----G
		"GTCTT", // -----
	})
	assertNormal(t, []string{
		"CTGTT", // -----
		"TGCGA", // -G---
		"GAGTA", // --G--
		"TAGCT", // ---!-
		"GACGG", // ----G
		"GTCTT", // -----
	})

	assertMutant(t, []string{
		"CTGTT", // -----
		"TGCGA", // -----
		"GTGTA", // ----A
		"TAGAT", // ---A-
		"GAAGG", // --A--
		"GACTT", // -A---
	})
	assertNormal(t, []string{
		"CTGTT", // -----
		"TGCGA", // -----
		"GTGTA", // ----A
		"TAGAT", // ---A-
		"GACGG", // --!--
		"GACTT", // -A---
	})
}
