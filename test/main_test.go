package test

import (
	"strings"
	"testing"

	"github.com/speedmancs/vmmanager/model"
)

func setup() {
	model.Clear()
	jsons := []string{
		`{ "name": "vm0", "status": "stopped" }`,
		`{ "name": "vm1", "status": "running" }`,
		`{ "name": "vm2", "status": "stopped" }`,
	}

	for _, json := range jsons {
		r := strings.NewReader(json)
		model.RegisterVM(r)
	}
}

func assertString(t *testing.T, message string, actual string, expected string) {
	if actual != expected {
		t.Errorf("%s: actual:%s, expected:%s", message, actual, expected)
	}
}

func assertInt(t *testing.T, message string, actual int, expected int) {
	if actual != expected {
		t.Errorf("%s: actual:%d, expected:%d", message, actual, expected)
	}
}

func TestRegisterVM(t *testing.T) {
	model.Clear()
	const jsonStream = `
	{ "name": "vm0", "status": "stopped" }
	`
	r := strings.NewReader(jsonStream)
	vm, _ := model.RegisterVM(r)
	assertInt(t, "vm.ID", vm.ID, 0)
	assertString(t, "vm.Name", vm.Name, "vm0")
	assertString(t, "vm.Status", vm.Status, "stopped")
}

func TestGetAllVMs(t *testing.T) {
	setup()
	vms, _ := model.GetVMs()
	assertInt(t, "vms count", len(vms), 3)
	assertInt(t, "vms[1].Id", vms[1].ID, 1)
	assertString(t, "vms[1].Name", vms[1].Name, "vm1")
	assertString(t, "vms[1].Status", vms[1].Status, "running")
}

func TestGetVM(t *testing.T) {
	setup()
	vm, _ := model.GetVM(2)
	assertInt(t, "vm id", vm.ID, 2)
	assertString(t, "vm name", vm.Name, "vm2")
	assertString(t, "vm status", vm.Status, "stopped")

	_, err := model.GetVM(3)
	if err == nil {
		t.Errorf("Should have error when getting vm with id 3")
	}
}

func TestDeleteVM(t *testing.T) {
	setup()
	model.DeleteVM(2)
	assertInt(t, "vm count", model.Count(), 2)
	err := model.DeleteVM(2)
	if err == nil {
		t.Errorf("Should have error when getting vm with id 2")
	}
}

func TestUpdateVM(t *testing.T) {
	setup()
	const jsonStream = `
	{ "name": "vm1_new", "status": "stopped" }
	`
	r := strings.NewReader(jsonStream)
	vm, _ := model.UpdateVM(1, r)
	assertInt(t, "vm.ID", vm.ID, 1)
	assertString(t, "vm.Name", vm.Name, "vm1_new")
	assertString(t, "vm.Status", vm.Status, "stopped")

	r = strings.NewReader(jsonStream)
	_, err := model.UpdateVM(3, r)
	if err == nil {
		t.Errorf("Should have error when getting vm with id 2")
	}
}
