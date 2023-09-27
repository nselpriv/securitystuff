package main

import (
	"fmt"
	"math/big"
	"math/rand"
)


type values struct {
	public_shared_base big.Int //g
	public_shared_prime big.Int //p
	BOB_public_key big.Int //g^x mod p
}


func main() {
	message:= "’2000’."

	vals:=&values{
		public_shared_base: *big.NewInt(666),
		public_shared_prime: *big.NewInt(6661),
		BOB_public_key: *big.NewInt(2227),
	}
	
	elgamal(vals, message)
}

func elgamal(v *values, s string){
	//first we select a random y for alice
	
	r := rand.Intn(6661 - 666) + 666
	bigR := big.NewInt(int64(r))
	fmt.Print(r)
	
	fmt.Println(new(big.Int).Exp(bigR, &v.BOB_public_key, &v.public_shared_base))
	
}