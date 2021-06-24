package test

import (
	"strings"
	"testing"

	"github.com/speedmancs/vmmanager/model"
	"github.com/stretchr/testify/assert"
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

func TestRegisterVM(t *testing.T) {
	model.Clear()
	const jsonStream = `
	{ "name": "vm0", "status": "stopped" }
	`
	r := strings.NewReader(jsonStream)
	vm, _ := model.RegisterVM(r)
	assert.Equal(t, 0, vm.ID, "VM ID should be 0")
	assert.Equal(t, "vm0", vm.Name, "VM Name should be vm0")
	assert.Equal(t, "stopped", vm.Status, "VM Status should be stopped")
}

func TestGetAllVMs(t *testing.T) {
	setup()
	vms, _ := model.GetVMs()
	assert.Equal(t, 3, len(vms), "VMs count should be 3")
	assert.Equal(t, 1, vms[1].ID, "ID of vms[1] should be 1")
	assert.Equal(t, "vm1", vms[1].Name, "Name of VM should be vm1")
	assert.Equal(t, "running", vms[1].Status, "Status of VM should be running")
}

func TestGetVM(t *testing.T) {
	setup()
	vm, _ := model.GetVM(2)
	assert.Equal(t, 2, vm.ID, "ID of vms[1] should be 2")
	assert.Equal(t, "vm2", vm.Name, "Name of VM should be vm2")
	assert.Equal(t, "stopped", vm.Status, "Status of VM should be stopped")

	_, err := model.GetVM(3)
	assert.NotNil(t, err, "Should have error when getting vm with id 3")
}

func TestDeleteVM(t *testing.T) {
	setup()
	model.DeleteVM(2)
	assert.Equal(t, 2, model.Count(), "VMs count should be 2")
	err := model.DeleteVM(2)
	assert.NotNil(t, err, "Should have error when deleting vm with id 2")
}

func TestUpdateVM(t *testing.T) {
	setup()
	const jsonStream = `
	{ "name": "vm1_new", "status": "stopped" }
	`
	r := strings.NewReader(jsonStream)
	vm, _ := model.UpdateVM(1, r)
	assert.Equal(t, 1, vm.ID, "ID of vms[1] should be 1")
	assert.Equal(t, "vm1_new", vm.Name, "Name of VM should be vm1_new")
	assert.Equal(t, "stopped", vm.Status, "Status of VM should be stopped")

	r = strings.NewReader(jsonStream)
	_, err := model.UpdateVM(3, r)
	assert.NotNil(t, err, "Should have error when updating vm with id 3")
}
