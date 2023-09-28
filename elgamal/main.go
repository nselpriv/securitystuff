package main

import (
	"fmt"
	"math/big"
	"math/rand"

)


type values struct {
	g, //g
	p, //p
	BOB_public_key,  //g^x mod p
	m big.Int
}


func main() {
	vals:=&values{
		g: *big.NewInt(666),
		p: *big.NewInt(6661),
		BOB_public_key: *big.NewInt(2227),
		m: *big.NewInt(2000),
	}
	a := rand.Intn((6661-1) - 1) + 1
	pkAlice := big.NewInt(int64(a))

	c1, c2 := encrypt(vals,pkAlice)

	m:= decrypt(vals, pkAlice, c1, c2)

	fmt.Printf("message is %v", m.Int64())
}

func encrypt(v *values, pkAlice *big.Int) (c1, c2 *big.Int){
	
	fmt.Printf("message to be encrypted is %v\n", v.m.Int64())
	c1 = new(big.Int).Exp(&v.g, pkAlice, &v.p)
	c2 = new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(&v.BOB_public_key, pkAlice, &v.p), &v.m),&v.p)

	fmt.Printf("generated cipher 1 with value %v\n ",c1.Int64())
	fmt.Printf("generated cipher 2 with value %v\n ",c2.Int64())
	return

}

func decrypt(v *values, pk, c1, c2 *big.Int) (m *big.Int) {

	bobKey := crackKey(v)

	uglyminus := new(big.Int).Sub(new(big.Int).Sub(&v.p, big.NewInt(1)), bobKey)
	m = new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(c1, uglyminus, &v.p),c2),&v.p)
	return 
}

func crackKey(v *values) (key *big.Int) {


	return
}