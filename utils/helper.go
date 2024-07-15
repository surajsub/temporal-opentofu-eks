package utils

import (
	"log"
)

var provisioner Provisioner

func GetProvisioner(provisioner Provisioner, prov string) (Provisioner, string, string) {

	var engine string
	var enginedir string

	log.Println("Default flag is Opentofu")
	provisioner = &TofuProvisioner{}
	engine = "tofu"
	enginedir = BASEOTDIRECTORY

	log.Println("the engine return from the call is ", engine)
	return provisioner, engine, enginedir
}
