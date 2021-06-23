package model

type VM struct {
	ID     int
	Name   string
	Status string
}

var vms []VM

func GetVMs() ([]VM, error) {
	return vms, nil
}
