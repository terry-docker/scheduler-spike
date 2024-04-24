package common

import "log"

// TODO Instead of logging VolumeName we want the volume id and to export the volume.
type TaskConfig struct {
	ID         string `json:"id"`
	Spec       string `json:"spec"`
	VolumeName string `json:"VolumeName"`
}

func Cmd(name string) func() {
	return func() {
		log.Printf("Executing task: %s", name)
	}
}
