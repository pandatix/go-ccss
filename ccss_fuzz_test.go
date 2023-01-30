package goccss

import "testing"

func FuzzParseVector(f *testing.F) {
	for _, tt := range testsParseVector {
		f.Add(tt.Vector)
	}

	f.Fuzz(func(t *testing.T, vector string) {
		ccss, err := ParseVector(vector)

		if err != nil {
			if ccss != nil {
				t.Fatal("Not supposed to get a CCSS when an error is returned")
			}
		} else {
			// This check works because CCSS has a predetermined order
			ccssVector := ccss.Vector()
			if vector != ccssVector {
				t.Fatalf("vetor differs at export: input is %s but output is %s", vector, ccssVector)
			}
		}
	})
}
