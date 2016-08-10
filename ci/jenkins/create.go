package main

import "fmt"

func (r CreateReq) create_jenkins_sources() error {
    if r.Name == "" {
        return fmt.Errorf("Missing jenkins instance Name")
    }
    return nil
}
