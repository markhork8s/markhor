package decrypt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestParseOrderingRec_Valid(t *testing.T) {
	t.Parallel()
	raw := []string{
		"apiVersion",
		"kind",
		"metadata/name",
		"metadata/namespace",
		"markhorParams/hierarchySeparator",
		"markhorParams/managedAnnotation",
		"markhorParams/order",
		"type",
		"data/sessionSecret",
		"stringData/another",
	}

	expectedResult := []Ordering{
		{Name: "apiVersion", Terms: []Ordering{}},
		{Name: "kind", Terms: []Ordering{}},
		{Name: "metadata", Terms: []Ordering{
			{Name: "name", Terms: []Ordering{}},
			{Name: "namespace", Terms: []Ordering{}},
		}},
		{Name: "markhorParams", Terms: []Ordering{
			{Name: "hierarchySeparator", Terms: []Ordering{}},
			{Name: "managedAnnotation", Terms: []Ordering{}},
			{Name: "order", Terms: []Ordering{}},
		}},
		{Name: "type", Terms: []Ordering{}},
		{Name: "data", Terms: []Ordering{
			{Name: "sessionSecret", Terms: []Ordering{}},
		}},
		{Name: "stringData", Terms: []Ordering{
			{Name: "another", Terms: []Ordering{}},
		}},
	}

	result, err := parseOrderingRec(raw, "/")
	if err != nil {
		t.Fatal("This error whould have not occcourred", err)
	}
	diff := cmp.Diff(expectedResult, result)
	if diff != "" {
		t.Errorf("%s", diff)
	}
}

func TestParseOrderingRec_Invalid(t *testing.T) {
	t.Parallel()
	inputs := [][]string{
		{
			"apiVersion",
			"kind",
			"metadata/name",
			"metadata/namespace",
			"markhorParams/hierarchySeparator",
			"markhorParams/managedAnnotation",
			"markhorParams/order",
			"markhorParams", //culprit
			"type",
		},
		{
			"apiVersion",
			"kind",
			"metadata/name",
			"metadata/namespace",
			"markhorParams/hierarchySeparator",
			"markhorParams/managedAnnotation",
			"markhorParams/order/nested", //culprit
			"markhorParams/order",
			"type",
		},
	}
	for i, input := range inputs {
		innerinput := input
		t.Run(fmt.Sprint("Input number ", i, " (starting form 0)"), func(t *testing.T) {
			t.Parallel()
			_, err := parseOrderingRec(innerinput, "/")
			if err == nil {
				t.Fatal("This operation should have failed")
			}
		})
	}
}

func TestParseOrdering_Invalid_Check_Error_Duplicate(t *testing.T) {
	t.Parallel()
	culprit := "markhorParams/hierarchySeparator"
	raw := []string{
		"apiVersion",
		"kind",
		"metadata/name",
		"metadata/namespace",
		"markhorParams/hierarchySeparator",
		"markhorParams/managedAnnotation",
		"markhorParams/order",
		culprit,
		"type",
	}
	_, err := parseOrdering(raw, "/")
	errorMesage := err.Error()
	if !strings.Contains(errorMesage, duplicatesErrorMessage) {
		t.Fatal("Error message should say it found a duplicate, but instead reads: ", err)
	}
	if !strings.Contains(errorMesage, culprit) {
		t.Fatal("Error message should contain the culprit item, but instead reads: ", err)
	}
}

func TestParseOrdering_Valid(t *testing.T) {
	t.Parallel()
	raw := []string{
		"apiVersion",
		"kind",
		"metadata.name",
		"metadata.namespace",
		"markhorParams.hierarchySeparator",
		"markhorParams.managedAnnotation",
		"markhorParams.order",
		"type",
		"data.sessionSecret",
		"stringData.another",
	}

	expectedResult := []Ordering{
		{Name: "apiVersion", Terms: []Ordering{}},
		{Name: "kind", Terms: []Ordering{}},
		{Name: "metadata", Terms: []Ordering{
			{Name: "name", Terms: []Ordering{}},
			{Name: "namespace", Terms: []Ordering{}},
		}},
		{Name: "markhorParams", Terms: []Ordering{
			{Name: "hierarchySeparator", Terms: []Ordering{}},
			{Name: "managedAnnotation", Terms: []Ordering{}},
			{Name: "order", Terms: []Ordering{}},
		}},
		{Name: "type", Terms: []Ordering{}},
		{Name: "data", Terms: []Ordering{
			{Name: "sessionSecret", Terms: []Ordering{}},
		}},
		{Name: "stringData", Terms: []Ordering{
			{Name: "another", Terms: []Ordering{}},
		}},
	}

	result, err := parseOrdering(raw, ".")
	if err != nil {
		t.Fatal("This error whould have not occcourred", err)
	}
	diff := cmp.Diff(expectedResult, result)
	if diff != "" {
		t.Errorf("%s", diff)
	}
}

func TestParseOrdering_Invalid(t *testing.T) {
	t.Parallel()
	inputs := [][]string{
		{
			"apiVersion",
			"kind",
			"metadata.name",
			"metadata.namespace",
			"markhorParams.hierarchySeparator",
			"markhorParams.managedAnnotation",
			"markhorParams.order",
			"markhorParams", //culprit
			"type",
		},
		{
			"apiVersion",
			"kind",
			"metadata", //culprit
			"metadata.name",
			"metadata.namespace",
			"markhorParams.hierarchySeparator",
			"markhorParams.managedAnnotation",
			"markhorParams.order",
			"type",
		},
	}
	for i, input := range inputs {
		innerinput := input
		t.Run(fmt.Sprint("Input number ", i, " (starting form 0)"), func(t *testing.T) {
			t.Parallel()
			_, err := parseOrdering(innerinput, ".")
			if err == nil {
				t.Fatal("This operation should have failed")
			}
		})
	}
}

func TestSortWithOrdering(t *testing.T) {
	t.Parallel()
	t.Run("Successful ordering with valid simple input", func(t *testing.T) {
		t.Parallel()
		jsonData := map[string]interface{}{
			"c": 3,
			"b": 1,
			"a": 2,
		}
		ordering := Ordering{
			Name: "",
			Terms: []Ordering{
				{Name: "a", Terms: nil},
				{Name: "b", Terms: nil},
				{Name: "c", Terms: nil},
			},
		}
		expected := orderedmap.New[string, interface{}]()
		expected.Set("a", 2)
		expected.Set("b", 1)
		expected.Set("c", 3)

		result, err := sortWithOrdering(jsonData, ordering)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Result doesn't match the expected ordering")
		}
	})

	t.Run("Successful ordering with valid complex input", func(t *testing.T) {
		t.Parallel()
		jsonData := map[string]interface{}{
			"apiVersion": "markhor.example.com/v1",
			"kind":       "MarkhorSecret",
			"metadata": map[string]interface{}{
				"namespace": "default",
				"name":      "sample-secret",
			},
			"markhorParams": map[string]interface{}{
				"order": []interface{}{
					"apiVersion",
					"kind",
					"metadata/name",
					"metadata/namespace",
					"markhorParams/hierarchySeparator",
					"markhorParams/managedAnnotation",
					"markhorParams/order",
					"type",
					"data/session_secret",
					"stringData/another",
				},
				"hierarchySeparator": "/",
				"managedAnnotation":  "",
			},
			"type": "Opaque",
			"stringData": map[string]interface{}{
				"another": "I want some pineapples",
			},
			"data": map[string]interface{}{
				"session_secret": "aHR0cHM6Ly95b3V0dS5iZS9kUXc0dzlXZ1hjUT8=",
			},
		}

		ordering := Ordering{
			Name: "",
			Terms: []Ordering{
				{Name: "apiVersion", Terms: []Ordering{}},
				{Name: "kind", Terms: []Ordering{}},
				{Name: "metadata", Terms: []Ordering{
					{Name: "name", Terms: []Ordering{}},
					{Name: "namespace", Terms: []Ordering{}},
				}},
				{Name: "markhorParams", Terms: []Ordering{
					{Name: "hierarchySeparator", Terms: []Ordering{}},
					{Name: "managedAnnotation", Terms: []Ordering{}},
					{Name: "order", Terms: []Ordering{}},
				}},
				{Name: "data", Terms: []Ordering{
					{Name: "session_secret", Terms: []Ordering{}},
				}},
				{Name: "stringData", Terms: []Ordering{
					{Name: "another", Terms: []Ordering{}},
				}},
				{Name: "type", Terms: []Ordering{}},
			},
		}

		expected := orderedmap.New[string, interface{}]()
		{
			expected.Set("apiVersion", "markhor.example.com/v1")
			expected.Set("kind", "MarkhorSecret")
			metadata := orderedmap.New[string, interface{}]()
			metadata.Set("name", "sample-secret")
			metadata.Set("namespace", "default")
			expected.Set("metadata", metadata)
			markhorParams := orderedmap.New[string, interface{}]()
			markhorParams.Set("hierarchySeparator", "/")
			markhorParams.Set("managedAnnotation", "")
			markhorParams.Set("order", []string{
				"apiVersion",
				"kind",
				"metadata/name",
				"metadata/namespace",
				"markhorParams/hierarchySeparator",
				"markhorParams/managedAnnotation",
				"markhorParams/order",
				"type",
				"data/session_secret",
				"stringData/another",
			})
			expected.Set("markhorParams", markhorParams)
			data := orderedmap.New[string, interface{}]()
			data.Set("session_secret", "aHR0cHM6Ly95b3V0dS5iZS9kUXc0dzlXZ1hjUT8=")
			expected.Set("data", data)
			stringData := orderedmap.New[string, interface{}]()
			stringData.Set("another", "I want some pineapples")
			expected.Set("stringData", stringData)
			expected.Set("type", "Opaque")
		}

		result, err := sortWithOrdering(jsonData, ordering)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		m1, _ := json.Marshal(result)
		m2, _ := json.Marshal(expected)
		s1 := string(m1)
		s2 := string(m2)
		if s1 != s2 {
			t.Errorf("Result doesn't match the expected ordering. Expected %s, got %s", s1, s2)
		}
	})

	t.Run("Unsuccessful ordering with invalid input", func(t *testing.T) {
		t.Parallel()
		jsonData := map[string]interface{}{
			"a": 2,
			"b": 1,
			"c": 3,
		}
		ordering := Ordering{
			Name: "",
			Terms: []Ordering{
				{Name: "a", Terms: nil},
				{Name: "b", Terms: nil},
				{Name: "d", Terms: nil}, // 'd' key doesn't exist in jsonData
			},
		}

		result, err := sortWithOrdering(jsonData, ordering)

		if err == nil {
			t.Error("Expected an error, but got none")
		}

		if result != nil {
			t.Error("Expected result to be nil")
		}
	})
}

// Sorts according to the order and puts sops at the end
func TestSortJson(t *testing.T) {
	t.Parallel()
	jsonData := map[string]interface{}{
		"apiVersion": "markhor.example.com/v1",
		"kind":       "MarkhorSecret",
		"metadata": map[string]interface{}{
			"namespace": "default",
			"name":      "sample-secret",
		},
		"markhorParams": map[string]interface{}{
			"order": []interface{}{
				"apiVersion",
				"kind",
				"metadata/name",
				"metadata/namespace",
				"markhorParams/hierarchySeparator",
				"markhorParams/managedAnnotation",
				"markhorParams/order",
				"type",
				"data/session_secret",
				"stringData/another",
			},
			"hierarchySeparator": "/",
			"managedAnnotation":  "",
		},
		"sops": map[string]interface{}{
			"something": "values",
		},
		"type": "Opaque",
		"stringData": map[string]interface{}{
			"another": "I want some pineapples",
		},
		"data": map[string]interface{}{
			"session_secret": "aHR0cHM6Ly95b3V0dS5iZS9kUXc0dzlXZ1hjUT8=",
		},
	}

	expected := orderedmap.New[string, interface{}]()
	{
		expected.Set("apiVersion", "markhor.example.com/v1")
		expected.Set("kind", "MarkhorSecret")
		metadata := orderedmap.New[string, interface{}]()
		metadata.Set("name", "sample-secret")
		metadata.Set("namespace", "default")
		expected.Set("metadata", metadata)
		markhorParams := orderedmap.New[string, interface{}]()
		markhorParams.Set("hierarchySeparator", "/")
		markhorParams.Set("managedAnnotation", "")
		markhorParams.Set("order", []string{
			"apiVersion",
			"kind",
			"metadata/name",
			"metadata/namespace",
			"markhorParams/hierarchySeparator",
			"markhorParams/managedAnnotation",
			"markhorParams/order",
			"type",
			"data/session_secret",
			"stringData/another",
		})
		expected.Set("markhorParams", markhorParams)
		data := orderedmap.New[string, interface{}]()
		expected.Set("type", "Opaque")
		data.Set("session_secret", "aHR0cHM6Ly95b3V0dS5iZS9kUXc0dzlXZ1hjUT8=")
		expected.Set("data", data)
		stringData := orderedmap.New[string, interface{}]()
		stringData.Set("another", "I want some pineapples")
		expected.Set("stringData", stringData)
		sops := orderedmap.New[string, interface{}]()
		sops.Set("something", "values")
		expected.Set("sops", sops)
	}

	result := sortJson(jsonData, slog.String("", ""))

	m1, _ := json.Marshal(result)
	m2, _ := json.Marshal(expected)
	s1 := string(m1)
	s2 := string(m2)
	if s1 != s2 {
		t.Errorf("Result doesn't match the expected ordering. Expected %s, got %s", s2, s1)
	}
}
