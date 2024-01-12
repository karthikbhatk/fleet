//go:build darwin && cgo

package keystore

/*
   #cgo LDFLAGS: -framework CoreFoundation -framework Security
   #include <CoreFoundation/CoreFoundation.h>
   #include <Security/Security.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const service = "com.fleetdm.fleetd.enroll.secret"

var serviceStringRef = stringToCFString(service)

func Exists() bool {
	return true
}

func Name() string {
	return "default keychain"
}

// AddSecret will add a secret to the keychain. This secret can be retrieved by this application without any user authorization.
func AddSecret(secret string) error {

	query := C.CFDictionaryCreateMutable(
		C.kCFAllocatorDefault,
		0,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks,
	)
	defer C.CFRelease(C.CFTypeRef(query))

	data := C.CFDataCreate(C.kCFAllocatorDefault, (*C.UInt8)(unsafe.Pointer(C.CString(secret))), C.CFIndex(len(secret)))
	defer C.CFRelease(C.CFTypeRef(data))

	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecClass), unsafe.Pointer(C.kSecClassGenericPassword))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecAttrService), unsafe.Pointer(serviceStringRef))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecValueData), unsafe.Pointer(data))

	status := C.SecItemAdd(C.CFDictionaryRef(query), nil)
	if status != C.errSecSuccess {
		return fmt.Errorf("failed to add %v to keychain: %v", service, status)
	}
	return nil
}

// UpdateSecret will update a secret in the keychain. This secret can be retrieved by this application without any user authorization.
func UpdateSecret(secret string) error {

	query := C.CFDictionaryCreateMutable(
		C.kCFAllocatorDefault,
		0,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks,
	)
	defer C.CFRelease(C.CFTypeRef(query))

	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecClass), unsafe.Pointer(C.kSecClassGenericPassword))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecAttrService), unsafe.Pointer(serviceStringRef))

	update := C.CFDictionaryCreateMutable(
		C.kCFAllocatorDefault,
		0,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks,
	)
	defer C.CFRelease(C.CFTypeRef(update))

	data := C.CFDataCreate(C.kCFAllocatorDefault, (*C.UInt8)(unsafe.Pointer(C.CString(secret))), C.CFIndex(len(secret)))
	defer C.CFRelease(C.CFTypeRef(data))
	C.CFDictionaryAddValue(update, unsafe.Pointer(C.kSecValueData), unsafe.Pointer(data))

	status := C.SecItemUpdate(C.CFDictionaryRef(query), C.CFDictionaryRef(update))
	if status != C.errSecSuccess {
		return fmt.Errorf("failed to update %v in keychain: %v", service, status)
	}
	return nil
}

// GetSecret will retrieve a secret from the keychain. If the secret was added by user or another application,
// then this application needs to be authorized to retrieve the secret.
func GetSecret() (string, error) {
	var query C.CFMutableDictionaryRef
	query = C.CFDictionaryCreateMutable(
		C.kCFAllocatorDefault,
		0,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks,
	)
	defer C.CFRelease(C.CFTypeRef(query))

	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecClass), unsafe.Pointer(C.kSecClassGenericPassword))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecReturnData), unsafe.Pointer(C.kCFBooleanTrue))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecMatchLimit), unsafe.Pointer(C.kSecMatchLimitOne))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecAttrLabel), unsafe.Pointer(serviceStringRef))

	var data C.CFTypeRef
	status := C.SecItemCopyMatching(C.CFDictionaryRef(query), &data)

	if status != C.errSecSuccess {
		if status == C.errSecItemNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to retrieve %v from keychain: %v", service, status)
	}

	secret := C.CFDataGetBytePtr(C.CFDataRef(data))
	return C.GoString((*C.char)(unsafe.Pointer(secret))), nil
}

// deleteSecret will delete a secret from the keychain.
// This function is only used by tests. It is here because usage of CGO in tests is not supported.
func deleteSecret() error {
	query := C.CFDictionaryCreateMutable(
		C.kCFAllocatorDefault,
		0,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks,
	)
	defer C.CFRelease(C.CFTypeRef(query))

	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecClass), unsafe.Pointer(C.kSecClassGenericPassword))
	C.CFDictionaryAddValue(query, unsafe.Pointer(C.kSecAttrService), unsafe.Pointer(serviceStringRef))

	status := C.SecItemDelete(C.CFDictionaryRef(query))
	if status != C.errSecSuccess {
		return fmt.Errorf("failed to delete %v from keychain: %v", service, status)
	}
	return nil
}

// stringToCFString will return a CFStringRef
func stringToCFString(s string) C.CFStringRef {
	bytes := []byte(s)
	ptr := (*C.UInt8)(&bytes[0])
	return C.CFStringCreateWithBytes(C.kCFAllocatorDefault, ptr, C.CFIndex(len(bytes)), C.kCFStringEncodingUTF8, C.false)
}
