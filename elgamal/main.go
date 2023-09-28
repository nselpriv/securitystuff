package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
)

//this implementation is written in go,
//I did not expect that the big.int structure would be so messed up when doing mul and mod.
//it structured in a way where 1+1 becomes x.add(1,1) so 1+1+1 is then x.add(x.add(1,1),x)
//which makes it kind of unreadable. I tried to make some good print statements to show that it works.
//run the code with no flag to see part 1 & 2 using go run ./ in the sec folder or go run ./ -f true to see part 3

type values struct {
	g, //g
	p, //p
	BOB_public_key,  //g^x mod p
	m big.Int //message
}

var(
	f = flag.Bool("f", false, "run with -f true to simulate the interception")
)

func main() {
	//I make a struct to contain the different values so i dont have to pass them around 

	flag.Parse()
	vals:=&values{
		g: *big.NewInt(666),
		p: *big.NewInt(6661),
		BOB_public_key: *big.NewInt(2227),
		m: *big.NewInt(2000),
	}
	//generate a private key for alice 
	a := rand.Intn((6661-1) - 1) + 1
	pkAlice := big.NewInt(int64(a))

	//This is part 1 
	c1, c2 := encrypt(vals,pkAlice)

	if (!*f) {
		//part 2 
		m:= decrypt(vals, pkAlice, c1, c2)
		fmt.Printf("message is %v\n", m.Int64())
	} else {
		//part 3 
		c1f, c2f := intercept(c1,c2)
		mf := decrypt(vals, pkAlice, c1f, c2f)
		fmt.Printf("fake message is %v\n", mf.Int64())
	}	
}

/**
Encrypt takes the following inputs

a pointer to the value struct
alice private key as a pointer to a big integer 
and returns 
ciphertext 1 and ciphertext 2 as big integers 
**/
func encrypt(v *values, pkAlice *big.Int) (c1, c2 *big.Int){
	
	fmt.Printf("message to be encrypted is %v\n", v.m.Int64())
	//c1 = g^pkAlice mod p // Alice public key
	c1 = new(big.Int).Exp(&v.g, pkAlice, &v.p)
	//c2 = bob public key ^ pkAlice mod p * m mod p //encrypted message 
	c2 = new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(&v.BOB_public_key, pkAlice, &v.p), &v.m),&v.p)

	fmt.Printf("generated cipher 1 with value %v\n ",c1.Int64())
	fmt.Printf("generated cipher 2 with value %v\n ",c2.Int64())
	return
}


/**
decrypt takes 

a pointer to the value struct 
alice private key 
the two ciphertexts from encryt

returns the message. 

**/
func decrypt(v *values, pk, c1, c2 *big.Int) (m *big.Int) {

	bobKey := crackKey(v)
	//uglyminus is just p-1-bobPK, big.Int.Sub just makes it really ugly 
	uglyminus := new(big.Int).Sub(new(big.Int).Sub(&v.p, big.NewInt(1)), bobKey)
	//m = c1 ^uglyminus mod p * c2 mod p 
	m = new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(c1, uglyminus, &v.p),c2),&v.p)
	return 
}

/**
Crackkey is using the values from my value struct and returning the private key of bob. 

**/
func crackKey(v *values) (key *big.Int) {
	//to crack the key i just check if g^key mod p is equal to bobs public key
	//starting from zero and going up. 
	//this is trivial because we have such a small prime 
	key = big.NewInt(0)
	for{
		if new(big.Int).Exp(&v.g, key, &v.p).Cmp(&v.BOB_public_key) == 0 {
			return
		} else {
			key.Add(key, big.NewInt(1))
		}
	}
}


/**
intercept is using the two ciphertexts 
and simply modifies c2 to multiply the number with 2
**/
func intercept(c1,c2 *big.Int) (c1Fake, c2Fake *big.Int){
	c1Fake =c1 // alice public key
	c2Fake =c2.Mul(c2,big.NewInt(2)) // message timed with two, this works because its just a number
	//in a situation where the message is a string this would break.
	return
} 