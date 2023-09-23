package pkg

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	png := Generate()
	err := os.WriteFile("ldb.test.png", png, 0644)
	if err != nil {
		t.Fatalf("error writing to file")
	}
	if len(png) <= 0 {
		t.Fatalf("generated image was empty")
	}
	fmt.Println("Output written to ldb.test.png")
}
