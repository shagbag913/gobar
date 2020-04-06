package main

import (
    "fmt"
    "os"
)

/* Prints error to stderr and returns true if error is not nil */
func logFatal(err error) bool {
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return true
    }
    return false
}
