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
   * Neither the name of Pass Desktop nor the names of its
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

package util

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
)

type Configuration struct {
    Encrypted struct {
        Token string `json:"token"`
        Salt  string `json:"salt"`
    } `json:"encrypted"`

    Host string `json:"host"`
    Port int    `json:"port"`
    CA   string `json:"ca"`
}

const filename = "/config/config.json"
const iconpath = "/iconpack/"

func GetConfig(path string) (Configuration, error) {
    cfg := path + filename

    if _, err := os.Stat(cfg); os.IsNotExist(err) {
        log.Println("No config file present.")
        return Configuration{}, err
    } else {
        // Load config.
        return LoadConfiguration(cfg), nil
    }
}

func LoadConfiguration(file string) Configuration {
    var config Configuration

    // Try to open configuration file.
    configFile, err := os.Open(file)
    defer configFile.Close()

    // Bail out if there was an error reading it.
    if err != nil {
        log.Fatal(err.Error())
    }

    // Decode the config.
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)

    return config
}

func ListAvailableIcons(path string) map[string]bool {
    files, err := ioutil.ReadDir(path + iconpath)

    if err != nil {
        log.Fatal(err)
    }

    icons := map[string]bool{}

    for _, f := range files {
        icon := f.Name()
        icons[strings.TrimSuffix(icon, filepath.Ext(icon))] = true
    }

    return icons
}
