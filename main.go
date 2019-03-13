package main 

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
) 

/* Runs two tests, one with code that will compile, one with code which will not */
func main() {
	_runTest1()
	fmt.Print("\n------------------------------------------------------\n")
	_runTest2()
}

/* Runs code that will compile */
func _runTest1() {
	fmt.Println("\nTest 1\nThe following code will compile:\n")

	// The program that will compile
	program := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Print(\"Hello World\")\n}"

	fmt.Println(program + "\n")

	// Runs the compileProgram function and captures the output and whether or not the output is an error
	message, err := compileProgram(program)

	// Prints out that an error occurred if output is error, this code will not run
	if err {
		fmt.Println("A compile error occurred")
	}

	// Prints the output
	fmt.Println(message)
}

func _runTest2() {
	fmt.Println("\nTest 2\nThe following code will not compile:\n")

	// The program that will not compile
	program := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Prnt(\"Hello World\")\n}"

	fmt.Println(program + "\n")

	// Runs the compileProgram function and captures the output and whether or not the output is an error
	message, err := compileProgram(program)

	// Prints out that an error occurred if output is error, this code will run
	if err {
		fmt.Println("A compile error occurred")
	}

	// Prints the error message
	fmt.Println(message)
}

func compileProgram(program string) (string, bool) {
	// Sets the endpoint and HTTP method
	url := "https://golang.org/compile"
	method := "POST"

	// Sets the request body using the program
	payload := strings.NewReader("---boundary\n" +
								 "Content-Disposition: form-data; name=\"body\"\n\n" +
								 program + "\n\n---boundary--")
	// Creates a reference to the HTTP client object
	client := &http.Client {
		CheckRedirect: func (req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// Creates request object
	req, err := http.NewRequest(method, url, payload)

	// Prints the error message if something goes wrong
	if err != nil {
		fmt.Println(err)
	}

	// Adds the header for the boundary in the request body
	req.Header.Add("content-type", "multipart/form-data; boundary=-boundary")

	// Does the request, closes the request, and extracts the body
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Maps the output of the body to the data map
	var data map[string]string
	json.Unmarshal(body, &data)

	// Creates and sets boolean variable to tell if there was a compiler error from input
	var comp_err bool
	if data["compile_errors"] != "" {
		comp_err = true
	} else {
		comp_err = false
	}

	// Returns either the compile error message or string and whether a compile error occurred
	return string(data["compile_errors"]) + string(data["output"]), comp_err
}
