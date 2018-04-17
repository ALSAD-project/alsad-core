// Example use:

// The following example starts a data source as nc in Shell A,
// and starts the driver in Shell B. The user program "python 
// sgdclassifier.py" will be deployed by the driver. The driver
// will wire up the data source and the user program. The user 
// program should initiate a TCP socket to this driver.

// (Shell A)
// $ nc -lk 9999

// (Shell B)
// $ export StreamInURL=":9999"
// $ export StreamOutURL=":8888"
// $ export UserProgram="python sgdclassifier.py"
// $ ./driver


package main

import (
    "bufio"
    "fmt"
    "strings"
    "net"
    "log"
    "os"
    "os/exec"
    "os/signal"
    "syscall"

    "github.com/kelseyhightower/envconfig"
)

func signalHandler(cmd *exec.Cmd) {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT)
    go func() {
        for {
            s := <-signalChan
            switch s {
            case syscall.SIGINT:
                if err := cmd.Process.Kill(); err != nil {
                    panic(err)
                }
                os.Exit(0)
            }
        }
    }()
}

func handleStreamIn(dataSourceConn net.Conn, userProgConn net.Conn) {
    for {
        streamIn, err := bufio.NewReader(dataSourceConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(userProgConn, streamIn)
    }
}

func handleStreamOut(dataSourceConn net.Conn, userProgConn net.Conn) {
    for {
        streamOut, err := bufio.NewReader(userProgConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(dataSourceConn, streamOut)
    }
}

func main() {

    // Load ENV configurations
    driverConfig := config{}
    if err := envconfig.Process("driver", &driverConfig); err != nil {
        log.Fatalf("Error on processing configuration: %s", err.Error())
        return
    }

    // Run user program
    UserProg := strings.Split(driverConfig.UserProgram, " ")
    cmd := exec.Command(UserProg[0], UserProg[1:]...)
    if err := cmd.Start(); err != nil {
        panic(err)
    }
    // Create signal handler for killing user program when ^C received
    go signalHandler(cmd)

    // Connect to stream in data source
    dataSourceConn, err := net.Dial("tcp", driverConfig.StreamInURL)
    if err != nil {
        panic(err)
    }
    
    // Listen for stream out data target (i.e. user program)
    ln, err := net.Listen("tcp", driverConfig.StreamOutURL)
    if err != nil {
        panic(err)
    }

    // Always wait for user program(s) to connect
    for {
        // Once connected by a user program,
        userProgConn, err := ln.Accept()
        if err != nil {
            panic(err)
        }
        // pass stream in data to user program, and
        go handleStreamIn(dataSourceConn, userProgConn)
        // pass stream out data to data source. 
        go handleStreamOut(dataSourceConn, userProgConn)
    }

}