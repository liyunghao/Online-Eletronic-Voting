package main

import (
	"encoding/base64"
	"fmt"
)

/*
#cgo LDFLAGS: -lsodium
#include <sodium.h>
*/
import "C"

func main() {

	// Initialize sodium
	if C.sodium_init() == -1 {
		fmt.Println("sodium_init() failed")
		return
	}

	// Generate Sodium Sign Key
	var signKey [C.crypto_sign_SECRETKEYBYTES]byte
	var signPublicKey [C.crypto_sign_PUBLICKEYBYTES]byte
	if C.crypto_sign_keypair(
		(*C.uchar)(&signPublicKey[0]),
		(*C.uchar)(&signKey[0]),
	) == -1 {
		fmt.Println("crypto_sign_keypair() failed")
		return
	}

	// Output key in base64
	fmt.Println("Sign Key:")
	fmt.Println(signKey)
	fmt.Println(base64.StdEncoding.EncodeToString(signKey[:]))
	fmt.Println("Sign Public Key:")
	fmt.Println(signPublicKey)
	fmt.Println(base64.StdEncoding.EncodeToString(signPublicKey[:]))

	mes := []byte("Hello LIBSodium")

	// Generate Sodium Sign Message in Detached Mode
	var signature [C.crypto_sign_BYTES]byte
	if C.crypto_sign_detached(
		(*C.uchar)(&signature[0]),
		nil,
		(*C.uchar)(&mes[0]),
		C.ulonglong(len(mes)),
		(*C.uchar)(&signKey[0]),
	) == -1 {
		fmt.Println("crypto_sign_detached() failed")
		return
	}

	// Verify Signature
	if C.crypto_sign_verify_detached(
		(*C.uchar)(&signature[0]),
		(*C.uchar)(&mes[0]),
		C.ulonglong(len(mes)),
		(*C.uchar)(&signPublicKey[0]),
	) == 0 {
		fmt.Println("Signature Verified")
	} else {
		fmt.Println("Signature Verification Failed")
	}
}
