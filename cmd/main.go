package main

import (
	"fmt"
	engine "github.com/tejas-techstack/huus/internal/engine"
	//"os"
)

func main() {

	fileName := "./example.db"
	tree, _ := engine.Open(fileName, 10, 4096)

	// for i := 0; i < 100; i++ {
	// 	err := tree.PutInt(i, i)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

 //  for i := 0; i < 100; i++ {
	// 	_, err := tree.DeleteInt(i)
	// 	if err != nil {
	// 		fmt.Println(err)
 //      return
	// 	}
	// }

	err := engine.BeginQueryLoop(tree)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	// if err := os.Remove(fileName); err != nil {
	// 	fmt.Println(err)
	// }

}
