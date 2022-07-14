package main

import (
	"log"
	"os"

	"github.com/j0rdan0/microsoft-related/ADE-scripts/core"
)

func main() {

	kv := new(core.KVData)
	core.GetDiskEncryptionType(kv)

	secret, err := core.GetSecret(kv)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	data, err := core.UnwrapSecret(secret, kv)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	core.WriteBEKFile(data)

}
