package main

import (
    "os"
    "log"
    "encoding/json"
)

type Config struct {
    Encrypted struct {
        Token string `json:"token"`
        Nonce string `json:"nonce"`
        Salt  string `json:"salt"`
    } `json:"encrypted"`
    Host string `json:"host"`
    Port string `json:"port"`
}

func LoadConfiguration(file string) Config {

    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()

    if err != nil {
        log.Println(err.Error())
    }

    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)

    return config
}

func Filter(vs []string, f func(string) bool) []string {
    vsf := make([]string, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

func Map(vs []string, f func(string) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}