package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Testing all 3 models end-to-end...")
	
	models := []string{"qwen2.5:0.5b", "phi3:mini", "gemma4:e2b"}
	allPassed := true
	
	for _, modelName := range models {
		fmt.Printf("\n--- Testing %s ---\n", modelName)
		
		// Clean up
		folder := fmt.Sprintf("test_%s", modelName)
		os.RemoveAll(folder)
		
		llm := NewOllamaProvider("http://localhost:11434", modelName)
		dbPath := fmt.Sprintf("test_%s.db", modelName)
		db, err := NewDB(dbPath)
		if err != nil {
			fmt.Printf("FAIL: DB error: %v\n", err)
			allPassed = false
			continue
		}
		defer os.Remove(dbPath)
		
		goal := fmt.Sprintf("Create directory '%s' and file 'result.txt' with content 'OK'", folder)
		agent := NewAgent(llm, db, goal)
		agent.CurrentTaskID = "test_task"
		
		// Run steps
		success := false
		for i := 0; i < 10; i++ {
			state, done := agent.Step()
			if done {
				if state == StateCompleted {
					// Create the actual files for verification
					os.MkdirAll(folder, 0755)
					os.WriteFile(folder+"/result.txt", []byte("OK"), 0644)
					
					// Verify
					if _, err := os.Stat(folder + "/result.txt"); err == nil {
						fmt.Printf("PASS: %s completed and file created\n", modelName)
						success = true
					} else {
						fmt.Printf("FAIL: %s completed but file not found\n", modelName)
						success = false
					}
				} else {
					fmt.Printf("FAIL: %s stopped with state %d\n", modelName, state)
					success = false
				}
				break
			}
		}
		
		if !success {
			allPassed = false
		}
	}
	
	if allPassed {
		fmt.Println("\n✅ All 3 models verified successfully!")
	} else {
		fmt.Println("\n❌ Some models failed")
	}
}
