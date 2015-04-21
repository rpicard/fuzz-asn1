package main

import (
    "crypto/rand"
    "log"
)

type EncodingRuleset interface {
    RandomBoolean() []byte
    //RandomInteger() []byte
    //RandomBitString() []byte
    //RandomOctetString() []byte
    //RandomNull() []byte
    //RandomObjectIdentifier() []byte
    //RandomReal() []byte
    //RandomEnumerated() []byte
}

type berEncoding struct {}

func (berEncoding) RandomBoolean() []byte {
    // we know it's going to be 3 bytes long
    result := make([]byte, 3)

    // first byte is tag
    result[0] = 0x01

    // second byte is length (of contents)
    result[1] = 0x01

    // third byte is contents
    // 
    // if false it is 0x00
    // if true it can have any non-zero value
    //
    // we will use a random byte
    randomByte := make([]byte, 1)
    _, err := rand.Read(randomByte)
    if err != nil {
        log.Fatal(err)
    }

    result[3] = randomByte[0]

    return result
}

func main() {

}


