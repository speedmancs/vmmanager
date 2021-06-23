package model

import (
	"encoding/json"
	"fmt"
	"io"
)

type VM struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var nextId = 0
var vms []VM

func GetVMs() ([]VM, error) {
	return vms, nil
}

func GetVM(id int) (*VM, error) {
	for i, vm := range vms {
		if vm.ID == id {
			return &vms[i], nil
		}
	}

	return nil, fmt.Errorf("Failed to Get VM with id %d", id)
}

func RegisterVM(r io.Reader) (*VM, error) {
	var newVM VM
	err := json.NewDecoder(r).Decode(&newVM)
	if err != nil {
		return nil, fmt.Errorf("Failed to register VM as there's parsing issue")
	}
	newVM.ID = nextId
	nextId = nextId + 1
	vms = append(vms, newVM)
	return &newVM, nil
}

func find(id int) int {
	idx := -1
	for i, vm := range vms {
		if vm.ID == id {
			idx = i
			break
		}
	}
	return idx
}

func UpdateVM(id int, r io.Reader) (*VM, error) {
	idx := find(id)
	if idx == -1 {
		return nil, fmt.Errorf("Failed to update VM with id %d", id)
	}

	var vm VM
	err := json.NewDecoder(r).Decode(&vm)
	if err != nil {
		return nil, fmt.Errorf("Failed to update VM with id %d as there's parsing issue", id)
	}

	vm.ID = id
	vms[idx] = vm
	return &vms[idx], nil
}

func DeleteVM(id int) error {
	idx := find(id)
	if idx == -1 {
		return fmt.Errorf("Failed to delete VM with id %d", id)
	}

	vms = append(vms[0:idx], vms[idx+1:]...)
	return nil
}
