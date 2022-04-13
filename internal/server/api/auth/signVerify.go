package auth

/*
#cgo LDFLAGS: -lsodium
#include <sodium.h>
*/
import "C"

func isValidSignature(publicKey, message, signature []byte) bool {
	return C.crypto_sign_verify_detached(
		(*C.uchar)(&signature[0]),
		(*C.uchar)(&message[0]),
		C.ulonglong(len(message)),
		(*C.uchar)(&publicKey[0]),
	) == 0
}
