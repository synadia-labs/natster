package main

import (
	"github.com/nats-io/nkeys"
)

func main() {}

//export getXKey
func getXKey() string {
	pair, _ := nkeys.CreateCurveKeys()
	pub, _ := pair.PublicKey()
	seed, _ := pair.Seed()

	ret := make(map[string]interface{})
	ret["public"] = pub
	ret["seed"] = seed

	println("msg from inside wasm")

	return "hi"
}
