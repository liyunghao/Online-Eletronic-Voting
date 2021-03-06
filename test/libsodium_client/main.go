package main

import (
	"encoding/base64"
	"flag"
	"fmt"
)

/*
#cgo LDFLAGS: -lsodium
#include <sodium.h>
*/
import "C"

var (
	toGenerateKey = flag.Bool("generate", false, "Generate Sodium ed25519 Sign Key")
	toVerify      = flag.Bool("verify", false, "Verify Sodium ed25519 Sign Message")
)

func main() {
	flag.Parse()

	// Initialize sodium
	if C.sodium_init() == -1 {
		fmt.Println("sodium_init() failed")
		return
	}

	// Generate Sodium ed25519 Sign Key
	if *toGenerateKey || flag.NArg() < 2 {
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
		fmt.Println(base64.StdEncoding.EncodeToString(sign_key[:]))
		fmt.Println("Sign Public Key:")
		fmt.Println(base64.StdEncoding.EncodeToString(sign_pk[:]))

		return
	}

	if *toVerify {
		// Args[0] -> Sign Public Key
		// Args[1] -> Original Message
		// Args[2] -> Sign Message
		sign_pk, _ := base64.StdEncoding.DecodeString(flag.Arg(0))
		original_mes, _ := base64.StdEncoding.DecodeString(flag.Arg(1))
		sign_mes, _ := base64.StdEncoding.DecodeString(flag.Arg(2))

		// verify Signature
		if C.crypto_sign_verify_detached(
			(*C.uchar)(&sign_mes[0]),
			(*C.uchar)(&original_mes[0]),
			C.ulonglong(len(original_mes)),
			(*C.uchar)(&sign_pk[0]),
		) == 0 {
			fmt.Println("Signature Verified")
		} else {
			fmt.Println("Signature Verification Failed")
		}
		return
	}

	// Args[0] -> Sign Key
	// Args[1] -> message

	sign_key, _ := base64.StdEncoding.DecodeString(flag.Arg(0))
	mes, _ := base64.StdEncoding.DecodeString(flag.Arg(1))

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

	fmt.Println("Sign Message:")
	fmt.Println(base64.StdEncoding.EncodeToString(sign_mes[:]))
}
