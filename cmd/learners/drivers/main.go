// Example use:

// The following example starts a data source as nc in Shell A,
// and starts the driver in Shell B. The user program "python 
// sgdclassifier.py" will be deployed by the driver. The driver
// will wire up the data source and the user program. The user 
// program should initiate a TCP socket to this driver.

// $ export DISPATCHER_LISTEN_URL=":9999"
// $ export USERPROG_LISTEN_URL=":8888"
// $ export USER_PROGRAM="python etc/local/sgdclassifier.py"
// $ ./drivers


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
                fmt.Println("User program killed.")
                os.Exit(0)
            }
        }
    }()
}

func handleStreamIn(sourceConn net.Conn, targetConn net.Conn) {
    for {
        streamIn, err := bufio.NewReader(sourceConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(targetConn, streamIn)
        fmt.Println("Total bytes passed to user program:", len(streamIn))
    }
}

func handleStreamOut(sourceConn net.Conn, targetConn net.Conn) {
    for {
        streamOut, err := bufio.NewReader(targetConn).ReadString('\n')
        if err != nil {
            panic(err)
        }
        fmt.Fprintf(sourceConn, streamOut)
        fmt.Println("Total bytes passed to dispatcher:", len(streamOut))
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

    // [For testing] Connect to dispatcher
    // dispatcherConn, err := net.Dial("tcp", driverConfig.DispatcherListenURL)
    // if err != nil {
    //     panic(err)
    // }

    // Listen for dispatcher
    dispatcherLn, err := net.Listen("tcp", driverConfig.DispatcherListenURL)
    if err != nil {
        panic(err)
    }
    
    // Listen for user program
    userProgLn, err := net.Listen("tcp", driverConfig.UserProgListenURL)
    if err != nil {
        panic(err)
    }

    // Wait for user program to connect to this driver
    userProgConn, err := userProgLn.Accept()
    if err != nil {
        panic(err)
    }
    fmt.Println("User program connected.")

    // Always wait for dispatcher to connect
    for {
        // Once connected by dispatcher,
        dispatcherConn, err := dispatcherLn.Accept()
        if err != nil {
            panic(err)
        }
        fmt.Println("Dispatcher connected.")

        // pass stream in data to user program, and
        go handleStreamIn(dispatcherConn, userProgConn)
        // pass stream out data to dispatcher. 
        go handleStreamOut(dispatcherConn, userProgConn)
    }

}