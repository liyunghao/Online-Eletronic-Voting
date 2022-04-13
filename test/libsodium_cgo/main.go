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

	// Generate Sodium ed25519 Sign Key
	var sign_key [C.crypto_sign_ed25519_SECRETKEYBYTES]byte
	var sign_pk [C.crypto_sign_ed25519_PUBLICKEYBYTES]byte
	if C.crypto_sign_ed25519_keypair(
		(*C.uchar)(&sign_pk[0]),
		(*C.uchar)(&sign_key[0]),
	) == -1 {
		fmt.Println("crypto_sign_ed25519_keypair() failed")
		return
	}

	// Output key in base64
	fmt.Println("Sign Key:")
	fmt.Println(sign_key)
	fmt.Println(base64.StdEncoding.EncodeToString(sign_key[:]))
	fmt.Println("Sign Public Key:")
	fmt.Println(sign_pk)
	fmt.Println(base64.StdEncoding.EncodeToString(sign_pk[:]))

	mes := []byte("Hello LIBSodium")

	// Generate Sodium ed25519 Sign Message in Detached Mode
	var sign_mes [C.crypto_sign_ed25519_BYTES]byte
	if C.crypto_sign_detached(
		(*C.uchar)(&sign_mes[0]),
		nil,
		(*C.uchar)(&mes[0]),
		C.ulonglong(len(mes)),
		(*C.uchar)(&sign_key[0]),
	) == -1 {
		fmt.Println("crypto_sign_detached() failed")
		return
	}

	// verify Signature
	if C.crypto_sign_verify_detached(
		(*C.uchar)(&sign_mes[0]),
		(*C.uchar)(&mes[0]),
		C.ulonglong(len(mes)),
		(*C.uchar)(&sign_pk[0]),
	) == 0 {
		fmt.Println("Signature Verified")
	} else {
		fmt.Println("Signature Verification Failed")
	}
}
