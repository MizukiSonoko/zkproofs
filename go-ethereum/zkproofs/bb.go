// Copyright 2017 ING Bank N.V.
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package zkproofs

/*
This file contains the implementation of the BB signature scheme proposed in the paper:
Short signatures without random oracle
Boneh and Boyen
Eurocrypt 2004
*/

import (
	"crypto/rand"
	"errors"
	"github.com/MizukiSonoko/zkproofs/go-ethereum/crypto/bn256"
	"math/big"
)

type keypair struct {
	Pubk  *bn256.G1
	Privk *big.Int
}

/*
keygen is responsible for the key generation.
*/
func keygen() (keypair, error) {
	var (
		kp  keypair
		e   error
		res bool
	)
	kp.Privk, e = rand.Int(rand.Reader, bn256.Order)
	if e != nil {
		return kp, e
	}
	kp.Pubk, res = new(bn256.G1).Unmarshal(new(bn256.G1).ScalarBaseMult(kp.Privk).Marshal())
	if res {
		return kp, errors.New("Could not compute scalar multiplication.")
	}
	return kp, e
}

/*
sign receives as input a message and a private key and outputs a digital signature.
*/
func sign(m *big.Int, Privk *big.Int) (*bn256.G2, error) {
	var (
		res       bool
		signature *bn256.G2
	)
	inv := ModInverse(Mod(Add(m, Privk), bn256.Order), bn256.Order)
	signature, res = new(bn256.G2).Unmarshal(new(bn256.G2).ScalarBaseMult(inv).Marshal())
	if res != false {
		return signature, nil
	} else {
		return nil, errors.New("Error while computing signature.")
	}
}

/*
verify receives as input the digital signature, the message and the public key. It outputs
true if and only if the signature is valid.
*/
func verify(signature *bn256.G2, m *big.Int, Pubk *bn256.G1) (bool, error) {
	// e(y.g^m, sig) = e(g1,g2)
	var (
		gm     *bn256.G1
		e, res bool
	)
	// g^m
	gm, e = new(bn256.G1).Unmarshal(new(bn256.G1).ScalarBaseMult(m).Marshal())
	// y.g^m
	gm = gm.Add(gm, Pubk)
	// e(y.g^m, sig)
	p1 := bn256.Pair(gm, signature)
	// e(g1,g2)
	g1 := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(1))
	g2 := new(bn256.G2).ScalarBaseMult(new(big.Int).SetInt64(1))
	p2 := bn256.Pair(g1, g2)
	// p1 == p2?
	p2 = p2.Neg(p2)
	p1 = p1.Add(p1, p2)
	res = p1.IsOne()
	if e != false {
		return res, nil
	}
	return false, errors.New("Error while computing signature.")
}
