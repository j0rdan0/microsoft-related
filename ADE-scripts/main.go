package main

import (
	"log"
	"os"

	"github.com/j0rdan0/microsoft-related/ADE-scripts/core"
)

func main() {

	kv := new(core.KVData)
	core.GetDiskEncryptionType(kv)

	if ok, err := core.SetAccessPolicy(kv); !ok {
		log.Fatal(err)
	}
	token, err := core.GetToken()
	handleError(err)

	secret, err := core.GetSecret(token, kv)
	handleError(err)

	data, err := core.UnwrapSecret(secret, token, kv)
	handleError(err)
	core.WriteBEKFile(data)

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
