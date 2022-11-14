package core

import (
	"fmt"
	"testing"
)

func TestAddZerosPrefixTableDriven(t *testing.T) {
	// Given
	var tests = []struct {
		value string
		count int
		want  string
	}{
		{"1", -1, "1"},
		{"1", 0, "1"},
		{"1", 1, "1"},
		{"0", 2, "00"},
		{"10", 1, "10"},
		{"10", 2, "10"},
		{"10", 3, "010"},
		{"1", 3, "001"},
		{"1", 4, "0001"},
		{"1", 5, "00001"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s with %d", tt.value, tt.count)
		t.Run(testname, func(t *testing.T) {
			// When
			got := addZerosPrefix(tt.value, tt.count)

			// Then
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParseParamTableDrivenRelais(t *testing.T) {
	// Given
	tic := StandardTicValue{}
	var tests = []struct {
		value                                                  string
		want1, want2, want3, want4, want5, want6, want7, want8 int8
	}{
		{"000", 0, 0, 0, 0, 0, 0, 0, 0},
		{"001", 1, 0, 0, 0, 0, 0, 0, 0},
		{"002", 0, 1, 0, 0, 0, 0, 0, 0},
		{"140", 0, 0, 1, 1, 0, 0, 0, 1},
		{"255", 1, 1, 1, 1, 1, 1, 1, 1},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.value)
		t.Run(testname, func(t *testing.T) {
			// When
			tic.ParseParam("RELAIS", []string{tt.value, "!"})

			// Then
			if tic.Relai1 != tt.want1 {
				t.Errorf("got %d, want %d", tic.Relai1, tt.want1)
			}
			if tic.Relai2 != tt.want2 {
				t.Errorf("got %d, want %d", tic.Relai2, tt.want2)
			}
			if tic.Relai3 != tt.want3 {
				t.Errorf("got %d, want %d", tic.Relai3, tt.want3)
			}
			if tic.Relai4 != tt.want4 {
				t.Errorf("got %d, want %d", tic.Relai4, tt.want4)
			}
			if tic.Relai5 != tt.want5 {
				t.Errorf("got %d, want %d", tic.Relai5, tt.want5)
			}
			if tic.Relai6 != tt.want6 {
				t.Errorf("got %d, want %d", tic.Relai6, tt.want6)
			}
			if tic.Relai7 != tt.want7 {
				t.Errorf("got %d, want %d", tic.Relai7, tt.want7)
			}
			if tic.Relai8 != tt.want8 {
				t.Errorf("got %d, want %d", tic.Relai8, tt.want8)
			}
		})
	}
}

func TestParseDateTableDrivenRelais(t *testing.T) {
	// Given
	tic := StandardTicValue{}
	var tests = []struct {
		value string
		want  int64
	}{
		{"H000101000000", 946681200},
		{"H221113153547", 1668350147},
		{"H221113153548", 1668350148},
		{"E221218174516", 1671378316},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.value)
		t.Run(testname, func(t *testing.T) {
			// When
			tic.parseDate(tt.value)

			// Then
			if tic.Date.Unix() != tt.want {
				t.Errorf("got %d, want %d", tic.Date.Unix(), tt.want)
			}
		})
	}
}

func TestParseParam(t *testing.T) {
	// Given
	var values = []struct {
		name  string
		value []string
	}{
		{"PJOURF+1", []string{"00008001", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "NONUTILE", "9"}},
		{"ADSC", []string{"XXXXXX", "7"}},
		{"VTIC", []string{"02", "J"}},
		{"DATE", []string{"H221113153547", "D"}},
		{"NGTF", []string{"BASE", "<"}},
		{"EASTF", []string{"040626660E", "-", "F"}},
		{"EASF01", []string{"040393601", "<"}},
		{"EASF02", []string{"000233059", "9"}},
		{"EASF03", []string{"000000000", "$"}},
		{"EASF04", []string{"000000000", "%"}},
		{"EASF05", []string{"000000000", "&"}},
		{"EASF06", []string{"000000000", "'"}},
		{"EASF07", []string{"000000000", "("}},
		{"EASF08", []string{"000000000", ")"}},
		{"EASF09", []string{"000000000", "*"}},
		{"EASF10", []string{"000000000", "\""}},
		{"EASD01", []string{"040626660", ">"}},
		{"EASD02", []string{"000000000", "!"}},
		{"EASD03", []string{"000000000", "\""}},
		{"EASD04", []string{"000000000", "#"}},
		{"IRMS1", []string{"007", "5"}},
		{"URMS1", []string{"239", "H"}},
		{"PREF", []string{"06", "E"}},
		{"PCOUP", []string{"06", "_"}},
		{"SINSTS", []string{"01700", "N"}},
		{"SMAXSN", []string{"H221113002750", "01750", "2"}},
		{"SMAXSN-1", []string{"H221112151524", "01750", "S"}},
		{"CCASN", []string{"H221113150000", "01421", "3"}},
		{"CCASN-1", []string{"H221113140000", "01430", "P"}},
		{"UMOY1", []string{"H221113153000", "236", ","}},
		{"STGE", []string{"00DA0001", "K"}},
		{"MSG1", []string{"PAS DE", "MESSAGE", "<"}},
		{"PRM", []string{"16140520874326", "2"}},
		{"RELAIS", []string{"001", "B"}},
		{"NTARF", []string{"01", "N"}},
		{"NJOURF", []string{"00", "&"}},
		{"NJOURF+1", []string{"00", "B"}},
	}
	tic := StandardTicValue{}

	// When
	for _, testValues := range values {
		tic.ParseParam(testValues.name, testValues.value)
	}

	// Then
	// if tic.Date.UnixMilli() != 1668350147 {
	// 	t.Errorf("Expected date millis %d but got %d", 1668350147, tic.Date.UnixMilli())
	// }
	if tic.Relai1 != 1 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai2 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai3 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai4 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai5 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai6 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai7 != 0 {
		t.Error("Relais 1 not good")
	}
	if tic.Relai8 != 0 {
		t.Error("Relais 1 not good")
	}
}
