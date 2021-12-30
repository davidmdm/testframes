package testframes

import (
	"fmt"
	"os"
	"testing"
)

var calls []string

func BrandedTestingFunc(name string) func(*testing.T) {
	return func(t *testing.T) {
		calls = append(calls, name)
	}
}

func TestMain(m *testing.M) {
	if code := m.Run(); code != 0 {
		os.Exit(code)
	}

	expectedCalls := []string{
		"Before",
		"BeforeEach",
		"rootTest1",
		"AfterEach",
		"BeforeEach",
		"rootTest2",
		"AfterEach",
		"frame1/Before",
		"BeforeEach",
		"frame1/BeforeEach",
		"frameTest1",
		"AfterEach",
		"frame1/AfterEach",
		"BeforeEach",
		"frame1/BeforeEach",
		"frameTest2",
		"AfterEach",
		"frame1/AfterEach",
	}

	if len(expectedCalls) != len(calls) {
		panic("wrong call length")
	}

	for i := range expectedCalls {
		if expectedCalls[i] != calls[i] {
			panic(fmt.Sprintf("expected index %d to have call %q but got %q", i, expectedCalls[i], calls[i]))
		}
	}
}

func TestFrames(t *testing.T) {

	f := New(t)

	f.Before(BrandedTestingFunc("Before"))
	f.BeforeEach(BrandedTestingFunc("BeforeEach"))
	f.AfterEach(BrandedTestingFunc("AfterEach"))

	f.Run("root test 1", BrandedTestingFunc("rootTest1"))
	f.Run("root test 2", BrandedTestingFunc("rootTest2"))

	f.NextFrame("Frame", func(f *Frame) {

		f.BeforeEach(BrandedTestingFunc("frame1/BeforeEach"))
		f.Before(BrandedTestingFunc("frame1/Before"))
		f.AfterEach(BrandedTestingFunc("frame1/AfterEach"))

		f.Run("frame test 1", BrandedTestingFunc("frameTest1"))
		f.Run("frame test 2", BrandedTestingFunc("frameTest2"))
	})

}
