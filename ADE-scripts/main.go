package main

import (
	"log"

	"github.com/j0rdan0/microsoft-related/ADE-scripts/core"
)

func main() {

	kv := new(core.KVData)
	core.GetDiskEncryptionType(kv)

	if ok, err := core.SetAccessPolicy(kv); !ok {
		log.Fatal(err)
	}
	token, err := core.GetToken()
	core.HandleError(err)

	secret, err := core.GetSecret(token, kv)
	core.HandleError(err)

	data, err := core.UnwrapSecret(secret, token, kv)
	core.HandleError(err)
	core.WriteBEKFile(data)

}
