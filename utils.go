package main

import (
    "strconv"
    "bufio"
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

/*
 * Open configuration file from either (in order):
 *  a) Environment variable $GOBAR_CONFIG
 *  b) $XDG_CONFIG_HOME/gobar/gobar.conf
 *  c) $HOME/.config/gobar/gobar.conf
 */
func openConfigFile() (*os.File, error) {
    var configPath string

    gobarConfigEnv := os.Getenv("GOBAR_CONFIG")
    XdgConfigEnv := os.Getenv("XDG_CONFIG_HOME")
    if gobarConfigEnv != "" {
        configPath = gobarConfigEnv
    } else if XdgConfigEnv != "" {
        configPath = XdgConfigEnv + "/gobar/gobar.conf"
    } else {
        configPath = os.Getenv("HOME") + "/.config/gobar/gobar.conf"
    }

    file, err := os.Open(configPath)
    return file, err
}

func formatConfString(str string) string {
    newString := make([]byte, 0)
    for i := range str {
        if str[i] >= '!' && str[i] <= '~' {
            newString = append(newString, str[i])
        }
    }
    return string(newString)
}

func getConfValue(flag, defaultValue string) string {
    file, err := openConfigFile()
    if logFatal(err) {
        return ""
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    key := make([]byte, 0)
    value := make([]byte, 0)
    for scanner.Scan() {
        line := formatConfString(scanner.Text())
        passedEquals := false
        for i := range line {
            if passedEquals == false {
                if line[i] == '=' {
                    passedEquals = true
                } else {
                    key = append(key, line[i])
                }
                continue
            }
            value = append(value, line[i])
        }

        if string(key) != flag {
            value = nil
            key = nil
            continue
        } else {
            break
        }
    }

    if string(value) != "" {
        return string(value)
    } else {
        return defaultValue
    }
}

func getConfBool(flag string, defaultValue bool) bool {
    return getConfValue(flag, strconv.FormatBool(defaultValue)) == "true"
}

func getConfInt(flag string, defaultValue int) int {
    valueInt, _ := strconv.Atoi(getConfValue(flag, strconv.Itoa(defaultValue)))
    return valueInt
}
