package rand

import (
  // "fmt"
  "crypto/rand"
  "encoding/base64"
)

const RememberTokenBytes = 32

func Bytes(n int) ([]byte,error) {
  b := make([]byte, n)
  _, err := rand.Read(b)
  if err != nil {
    return nil, err
  }
  return b, nil
}

func String(nBytes int) (string, error) {
  b, err := Bytes(nBytes)
  if err != nil {
    return "", err
  }
  // fmt.Println("BYTE_SLICE:",b)
  str := base64.URLEncoding.EncodeToString(b)
  // fmt.Println("BASE64STRIND:",str)
  return str, nil
}

func RememberToken() (string,error) {
  return String(RememberTokenBytes)
}
