package slice

import (
	"testing"
)

func TestContains(t *testing.T) {
	// Test case 1: Check if a slice containing the element returns true
	slice1 := []int{1, 2, 3, 4, 5}
	element1 := 3
	if !Contains(slice1, element1) {
		t.Errorf("Expected slice %v to contain element %v", slice1, element1)
	}

	// Test case 2: Check if a slice not containing the element returns false
	slice2 := []string{"apple", "banana", "orange"}
	element2 := "grape"
	if Contains(slice2, element2) {
		t.Errorf("Expected slice %v not to contain element %v", slice2, element2)
	}

	// Test case 3: Check if an empty slice returns false
	var slice3 []int
	element3 := 0
	if Contains(slice3, element3) {
		t.Errorf("Expected empty slice %v to not contain any element", slice3)
	}

	// Test case 4: Check if a slice of structs containing a struct element returns true
	type person struct {
		Name string
		Age  int
	}
	slice4 := []person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 40},
		{Name: "Charlie", Age: 50},
	}
	element4 := person{Name: "Bob", Age: 40}
	if !Contains(slice4, element4) {
		t.Errorf("Expected slice %v to contain element %v", slice4, element4)
	}
}
