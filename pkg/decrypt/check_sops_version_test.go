package decrypt

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected string
		ok       bool
	}{
		{
			name:     "Valid case",
			input:    map[string]interface{}{"sops": map[string]interface{}{"version": "3.8"}},
			expected: "3.8",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, ok := getVersion(tc.input)
			if result != tc.expected || ok != true {
				t.Fatalf("Expected result: %v, Expected OK value: %v, Got result: %v, Got OK value: %v", tc.expected, tc.ok, result, ok)
			}
		})
	}

	invalidTestCases := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Missing sops field",
			input: map[string]interface{}{"sus": map[string]interface{}{"version": "3.8"}},
		},
		{
			name:  "Sops field is not an object",
			input: map[string]interface{}{"sops": "3.8.0"},
		},
		{
			name:  "Missing version field",
			input: map[string]interface{}{"sops": map[string]interface{}{"anotherversion": "3.8"}},
		},
		{
			name:  "Version field is not a string",
			input: map[string]interface{}{"sops": map[string]interface{}{"version": 38}},
		},
	}

	for _, testCase := range invalidTestCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, ok := getVersion(tc.input)
			if ok == true {
				t.Fatal("Expected GetVersion to fail with input ", tc.input)
			}
		})
	}
}

func TestIsSupportedSopsVersion(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		version  string
		expected string
		ok       bool
	}{
		{
			name:     "Supported version",
			version:  "3.8",
			expected: "",
			ok:       true,
		},
		{
			name:     "Supported version",
			version:  "3.8.1",
			expected: "",
			ok:       true,
		},
		{
			name:     "Supported version",
			version:  "3.8.1-alpha1",
			expected: "",
			ok:       true,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, ok := IsSupportedSopsVersion(tc.version)
			if result != tc.expected || ok != tc.ok {
				t.Fatalf("Expected result: %v, Expected OK value: %v, Got result: %v, Got OK value: %v", tc.expected, tc.ok, result, ok)
			}
		})
	}

	invalidTestCases := []struct {
		name    string
		version string
	}{
		{
			name:    "Invalid version: incorrect format",
			version: "3-8",
		},
		{
			name:    "Invalid version: unknown version",
			version: "-3.8",
		},
	}

	for _, testCase := range invalidTestCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, ok := IsSupportedSopsVersion(tc.version)
			if ok == true {
				t.Fatal("Expected IsSupportedSopsVersion to fail with version ", tc.version)
			}
		})
	}
}
