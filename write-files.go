package main

import (
	"fmt"
	"os"
)

func main() {
	// Use buffered writerds when writing lareg amounts of data (buf.io)
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Errorcreatingfile", file)
		return
	}
	defer file.Close()

	//write data to file
	data := []byte("Hwello World\n\n")
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}

	fmt.Println("Data has been  written to file successfully.")

	file, err = os.Create("WriteString.txt")

	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Dobry Den!\n")
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}

	fmt.Println("String has been written to file WriteString.txt successfully.")
}
