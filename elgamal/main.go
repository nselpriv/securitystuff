package main

import (
	"fmt"
	"math/big"
)


type values struct {
	public_shared_base big.Int
	public_shared_prime big.Int
	BOB_public_key big.Int
}


func main() {
	message:= "’2000’."

	vals:=values{
		public_shared_base: *big.NewInt(666),
		public_shared_prime: *big.NewInt(6661),
		BOB_public_key: *big.NewInt(2227),
	}
	



	fmt.Println(message)
	fmt.Println(vals.public_shared_base )
}