package main

import (
    "fmt"
    "log"
)

func testBTree() {
    fmt.Println("\n=== Starting B-tree Test Suite ===")
    
    // Test Case 1: Basic Insertion and Node Splitting
    fmt.Println("\nTest Case 1: Basic Insertion and Node Splitting")
    tree1 := createNewTree(3)
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    fmt.Println("Inserting numbers 1-10 sequentially")
    for _, num := range numbers {
        if err := tree1.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree1.PrintTree(tree1.root, 0)

    // Test Case 2: Duplicate Keys
    fmt.Println("\nTest Case 2: Duplicate Key Handling")
    tree2 := createNewTree(3)
    duplicates := []int{5, 5, 5, 5}
    fmt.Println("Attempting to insert duplicate keys")
    for _, num := range duplicates {
        if err := tree2.Insert(num); err != nil {
            log.Printf("Expected error for duplicate key %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree2.PrintTree(tree2.root, 0)

    // Test Case 3: Descending Order Insertion
    fmt.Println("\nTest Case 3: Descending Order Insertion")
    tree3 := createNewTree(3)
    descending := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
    fmt.Println("Inserting numbers 10-1 in descending order")
    for _, num := range descending {
        if err := tree3.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree3.PrintTree(tree3.root, 0)

    // Test Case 4: Large Number of Insertions
    fmt.Println("\nTest Case 4: Large Number of Insertions")
    tree4 := createNewTree(3)
    fmt.Println("Inserting 100 numbers")
    for i := 1; i <= 100; i++ {
        if err := tree4.Insert(i); err != nil {
            log.Printf("Error inserting %d: %v", i, err)
        }
    }
    fmt.Println("Tree height after 100 insertions:")
    tree4.PrintTree(tree4.root, 0)
    calculateAndPrintTreeHeight(tree4.root)

    // Test Case 5: Random Insertion Order
    fmt.Println("\nTest Case 5: Random Insertion Order")
    tree5 := createNewTree(3)
    random := []int{42, 15, 88, 3, 67, 91, 28, 53, 9, 33}
    fmt.Println("Inserting numbers in random order")
    for _, num := range random {
        if err := tree5.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree5.PrintTree(tree5.root, 0)

    // Test Case 6: Minimum Degree Edge Case
    fmt.Println("\nTest Case 6: Minimum Degree Edge Case")
    tree6 := createNewTree(2) // Minimum possible degree for B-tree
    minDegreeTest := []int{1, 2, 3, 4, 5}
    fmt.Println("Testing with minimum degree (t=2)")
    for _, num := range minDegreeTest {
        if err := tree6.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree6.PrintTree(tree6.root, 0)

    // Test Case 7: Large Minimum Degree
    fmt.Println("\nTest Case 7: Large Minimum Degree")
    tree7 := createNewTree(5) // Larger degree for wider nodes
    largeDegreeTest := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    fmt.Println("Testing with large degree (t=5)")
    for _, num := range largeDegreeTest {
        if err := tree7.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Resulting tree:")
    tree7.PrintTree(tree7.root, 0)
    
    // Test Case 8: Property Verification
    fmt.Println("\nTest Case 8: B-tree Property Verification")
    tree8 := createNewTree(3)
    propertyTest := []int{50, 25, 75, 12, 37, 62, 87, 6, 18, 31, 43, 56, 68, 81, 93}
    for _, num := range propertyTest {
        if err := tree8.Insert(num); err != nil {
            log.Printf("Error inserting %d: %v", num, err)
        }
    }
    fmt.Println("Verifying B-tree properties:")
    verifyBTreeProperties(tree8.root, tree8.minNode, nil, nil)
    
    fmt.Println("\n=== Test Suite Completed ===")
}

// Helper function to calculate and print tree height
func calculateAndPrintTreeHeight(node *Node) int {
    if node == nil {
        return 0
    }
    maxHeight := 0
    for _, child := range node.children {
        height := calculateAndPrintTreeHeight(child)
        if height > maxHeight {
            maxHeight = height
        }
    }
    return maxHeight + 1
}

// Helper function to verify B-tree properties
func verifyBTreeProperties(node *Node, t int, min, max *int) bool {
    // Check node key count
    if len(node.keys) > 2*t-1 {
        fmt.Printf("Node has too many keys: %d\n", len(node.keys))
        return false
    }
    
    if !node.isLeaf && len(node.keys) < t-1 {
        fmt.Printf("Internal node has too few keys: %d\n", len(node.keys))
        return false
    }

    // Check key ordering
    for i := 1; i < len(node.keys); i++ {
        if node.keys[i] <= node.keys[i-1] {
            fmt.Printf("Keys not in order: %d <= %d\n", node.keys[i], node.keys[i-1])
            return false
        }
    }

    // Check key range
    if min != nil && node.keys[0] <= *min {
        fmt.Printf("Key %d violates minimum bound %d\n", node.keys[0], *min)
        return false
    }
    if max != nil && node.keys[len(node.keys)-1] >= *max {
        fmt.Printf("Key %d violates maximum bound %d\n", node.keys[len(node.keys)-1], *max)
        return false
    }

    // Recursively check children
    if !node.isLeaf {
        for i, child := range node.children {
            var childMin, childMax *int
            if i > 0 {
                childMin = &node.keys[i-1]
            } else {
                childMin = min
            }
            if i < len(node.keys) {
                childMax = &node.keys[i]
            } else {
                childMax = max
            }
            if !verifyBTreeProperties(child, t, childMin, childMax) {
                return false
            }
        }
    }
    return true
}

