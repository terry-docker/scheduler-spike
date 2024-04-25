package common

import "log"

// TODO Is VolumeName unique? Should we be storing a volume id?
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
