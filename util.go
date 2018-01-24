/*
Copyright (c) 2018 Carl Löndahl. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
    "os"
    "io"
    "log"
    "encoding/json"
)

type Configuration struct {
    Encrypted struct {
        Token string `json:"token"`
        Nonce string `json:"nonce"`
        Salt  string `json:"salt"`
    } `json:"encrypted"`
    Host string `json:"host"`
    Port string `json:"port"`
}


func CopyFile(source string, destination string) error {
    s, err := os.Open(source)
    if err != nil {
        return err
    }

    // no need to check errors on read only file, we already got everything
    // we need from the filesystem, so nothing can go wrong now.
    defer s.Close()
    d, err := os.Create(destination)

    if err != nil {
        return err
    }

    if _, err := io.Copy(d, s); err != nil {
        d.Close()
        return err
    }

    return d.Close()
}

func LoadConfiguration(file string) Configuration {
    var config Configuration
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
