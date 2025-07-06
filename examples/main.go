package main

import (
	"fmt"
	"log"

	"github.com/demdxx/gocast/v2"
)

func main() {
	// Basic type casting examples
	fmt.Println("=== Basic Type Casting ===")
	fmt.Printf("String to int: %d\n", gocast.Number[int]("42"))
	fmt.Printf("Int to string: %s\n", gocast.Str(123))
	fmt.Printf("Float to int: %d\n", gocast.Number[int](3.14))

	// Deep copy examples
	fmt.Println("\n=== Deep Copy Examples ===")

	// 1. Simple struct copy
	type Person struct {
		Name string
		Age  int
	}

	original := Person{Name: "Alice", Age: 30}
	copied, err := gocast.TryCopy(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original: %+v\n", original)
	fmt.Printf("Copied: %+v\n", copied)

	// 2. Slice copy with specialized function
	originalSlice := []int{1, 2, 3, 4, 5}
	copiedSlice := gocast.CopySlice(originalSlice)

	// Modify original to show they're independent
	originalSlice[0] = 999
	fmt.Printf("Original slice: %v\n", originalSlice)
	fmt.Printf("Copied slice: %v\n", copiedSlice)

	// 3. Map copy with specialized function
	originalMap := map[string]int{"a": 1, "b": 2, "c": 3}
	copiedMap := gocast.CopyMap(originalMap)

	originalMap["a"] = 999
	fmt.Printf("Original map: %v\n", originalMap)
	fmt.Printf("Copied map: %v\n", copiedMap)

	// 4. Complex nested structure
	type Address struct {
		City    string
		Country string
	}

	type User struct {
		Name     string
		Age      int
		Address  Address
		Tags     []string
		Metadata map[string]any
	}

	complexUser := User{
		Name: "Bob",
		Age:  25,
		Address: Address{
			City:    "New York",
			Country: "USA",
		},
		Tags:     []string{"admin", "user"},
		Metadata: map[string]any{"role": "admin", "active": true},
	}

	copiedUser, err := gocast.TryCopy(complexUser)
	if err != nil {
		log.Fatal(err)
	}

	// Modify original to demonstrate deep copy
	complexUser.Address.City = "Los Angeles"
	complexUser.Tags[0] = "superadmin"
	complexUser.Metadata["role"] = "superadmin"

	fmt.Printf("Original user: %+v\n", complexUser)
	fmt.Printf("Copied user: %+v\n", copiedUser)

	// 5. Copy with options
	fmt.Println("\n=== Copy with Options ===")

	type Config struct {
		PublicSetting  string
		privateSetting string // unexported
		NestedConfig   *Config
	}

	config := &Config{
		PublicSetting:  "public",
		privateSetting: "private",
		NestedConfig: &Config{
			PublicSetting:  "nested_public",
			privateSetting: "nested_private",
		},
	}

	// Copy with options: ignore unexported fields, limit depth
	opts := gocast.CopyOptions{
		IgnoreUnexportedFields: true,
		MaxDepth:               1,
	}

	limitedCopy, err := gocast.TryCopyWithOptions(config, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original config: %+v\n", config)
	fmt.Printf("Limited copy: %+v\n", limitedCopy)
	fmt.Printf("Nested config copied: %v\n", limitedCopy.NestedConfig != nil)

	// 6. Circular reference handling
	fmt.Println("\n=== Circular Reference Handling ===")

	type Node struct {
		Value int
		Next  *Node
	}

	node1 := &Node{Value: 1}
	node2 := &Node{Value: 2}
	node1.Next = node2
	node2.Next = node1 // Create circular reference

	copiedNode, err := gocast.TryCopy(node1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original node1 value: %d\n", node1.Value)
	fmt.Printf("Copied node1 value: %d\n", copiedNode.Value)
	fmt.Printf("Circular reference preserved: %v\n", copiedNode.Next.Next == copiedNode)

	// 7. Performance comparison
	fmt.Println("\n=== Performance Examples ===")

	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	// Using general TryCopy
	_, err = gocast.TryCopy(slice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TryCopy: ✓")

	// Using specialized CopySlice (faster)
	_ = gocast.CopySlice(slice)
	fmt.Println("CopySlice: ✓ (faster)")

	// Using MustCopy (panics on error)
	_ = gocast.MustCopy(42)
	fmt.Println("MustCopy: ✓")

	fmt.Println("\n=== Examples completed successfully! ===")
}
