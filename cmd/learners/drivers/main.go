package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
)

func main() {
    
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    fmt.Println(text)

    cmd := exec.Command("bash", "-c", "./hello_world.sh")
    cmdOut, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(cmdOut))
}