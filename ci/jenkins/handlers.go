// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
    "fmt"
    "net/http"
    "github.hpe.com/christophe-larsonneur/goforjj"
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
)

// PluginData response object creator
func newPluginData() (* goforjj.PluginData) {
    var r goforjj.PluginData = goforjj.PluginData{
        Repos: make(map[string]goforjj.PluginRepo),
        Services: goforjj.PluginService{make(map[string]string)},
    }
    return &r
}

// Function to detect header content-type matching
// return true if match
func content_type_match(header http.Header, match string) bool {
    for _, v := range header["Content-Type"] {
        if (v == match) {
            return true
        }
    }
    return false
}

func panicIfError(w http.ResponseWriter, err error, message string, pars ...interface{}) {
    if err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if message != "" {
            err = fmt.Errorf("%s %s", fmt.Errorf(message, pars...), err)
        }
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
}


// Create handler
func Create(w http.ResponseWriter, r *http.Request) {
    var data *goforjj.PluginData = newPluginData()
    var req_data CreateReq

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))

    if err != nil {
        panic(err)
    }

    if content_type_match(r.Header, "application/json") {
        err := json.Unmarshal(body, &req_data)
        panicIfError(w, err, "Unable to decode '%#v' as json.", string(body))
    } else {
        panicIfError(w, *new(error), "Invalid payload format. Must be 'application/json'. Got %#v", r.Header["Content-Type"])
    }

    // Create the github.yaml source file.
    // See goforjj/plugin-json-struct.go for json data structure recognized by forjj.

    err_code := DoCreate(w, r, &req_data, data)

    req_data.SaveMaintainOptions(data)

    if data.ErrorMessage != "" {
        if err_code == 0 {
            err_code = 422 // unprocessable entity
        }
        log.Print("HTTP ERROR: ", err_code, " - ", data.ErrorMessage)
    } else {
        err_code = 200
    }
    w.WriteHeader(err_code)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
    }
}

// Update handler
func Update(w http.ResponseWriter, r *http.Request) {
    var data *goforjj.PluginData = newPluginData()
    var req_data UpdateReq

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))

    if err != nil {
        panic(err)
    }

    if content_type_match(r.Header, "application/json") {
        err := json.Unmarshal(body, &req_data)
        panicIfError(w, err, "Unable to decode '%#v' as json.", string(body))
    } else {
        panicIfError(w, *new(error), "Invalid payload format. Must be 'application/json'. Got %#v", r.Header["Content-Type"])
    }

    // Update the github.yaml source file.
    // See goforjj/plugin-json-struct.go for json data structure recognized by forjj.

    err_code := DoUpdate(w, r, &req_data, data)

    req_data.SaveMaintainOptions(data)

    if data.ErrorMessage != "" {
        if err_code == 0 {
            err_code = 422 // unprocessable entity
        }
        log.Print("HTTP ERROR: ", err_code, " - ", data.ErrorMessage)
    } else {
        err_code = 200
    }
    w.WriteHeader(err_code)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
    }
}

// Maintain handler
func Maintain(w http.ResponseWriter, r *http.Request) {
    var req_data MaintainReq
    var data *goforjj.PluginData = newPluginData()

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240))

    if err != nil {
        panic(err)
    }

    if content_type_match(r.Header, "application/json") {
        err := json.Unmarshal(body, &req_data)
        panicIfError(w, err, "Unable to decode '%#v' as json.", string(body))
    } else {
        panicIfError(w, *new(error), "Invalid payload format. Must be 'application/json'. Got %#v", r.Header["Content-Type"])
    }

    err_code := DoMaintain(w, r, &req_data, data)

    if data.ErrorMessage != "" {
        if err_code == 0 {
            err_code = 422 // unprocessable entity
        }
        log.Print("HTTP ERROR: ", err_code, " - ", data.ErrorMessage)
    } else {
        err_code = 200
    }
    w.WriteHeader(err_code)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
    }
}

// Index Handler
//
func Index(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprintf(w ,"FORJJ - jenkins driver for FORJJ. It is Implemented as a REST API.")
}

// Quit
func Quit(w http.ResponseWriter, _ *http.Request) {
    goforjj.DefaultQuit(w, "")
}
