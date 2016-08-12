package main

import (
    "github.hpe.com/christophe-larsonneur/goforjj"
    "log"
    "os"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) check_source_existence(ret *goforjj.PluginData) (status bool) {
    log.Printf("Checking Jenkins source code path existence.")

    if _, err := os.Stat(r.ForjjSourceMount) ; err == nil {
        ret.Errorf("Unable to maintain jenkins instances. '%s' is inexistent or innacessible.\n", r.ForjjSourceMount)
        return
    }
    ret.StatusAdd("environment checked.")
    return true
}

func (r *MaintainReq)instantiate(ret *goforjj.PluginData) (status bool) {

    return true
}
