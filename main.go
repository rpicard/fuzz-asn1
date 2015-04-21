package main

import (
    "crypto/rand"
    "log"
    "math/big"
)

type EncodingRuleset interface {
    RandomBoolean() []byte
    RandomInteger() []byte
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

func (berEncoding) RandomInteger() []byte {

    // we will choose a random length
    // then we will generate a random integer of that length
    //
    // our length will be between 0 and 127; we have 7 bits to specify it
    // using only one byte. the high bit of the length byte indiciates
    // whether or not more bytes are needed to specify length. if they are,
    // the rest of the length byte specifies how many additional bytes
    //
    // this means the real max length is the biggest integer than can be
    // specified with 128 bytes:
    //
    // 2^(8 * 128) = ~ 1.7 x 10^308
    //
    // in order to simplify for now, rather than randomly choosing how many
    // bytes to use for lenght AND randomly choosing a length that can be
    // represented by that many bytes AND randomly choosing a value that
    // can be represented by that length of bytes, we're just assuming
    // a single length byte
    //
    // TODO: ^ do that thing from the last paragraph

    // get a random int in [0, 128)
    bigLength, err := rand.Int(rand.Reader, big.NewInt(128))
    if err != nil {
        log.Fatal(err)
    }

    // convert to a regular old int
    length := int(bigLength.Int64())

    // generate our random content
    //
    // there might be a small issue with the rule to ensure that an integer
    // value is always encoded in the smallest possible number of octets
    //
    // http://www.itu.int/ITU-T/studygroups/com17/languages/X.690-0207.pdf
    //
    // 8.3.2:
    //
    // If the contents octets of an integer value encoding consist of more
    // than one octet, then the bits of the first octet and bit 8 of the
    // second octet:
    //
    //   a) shall not all be ones; and
    //   b) shall not all be zero.
    //
    // NOTE - These rules ensure that an integer value is always encoded in
    // the smallest possible number of octets
    //
    // the takeaway here is that we might be generating some technically
    // invalid BER encodings here. not sure what the consequences of that are
    //
    // TODO: does this need to be fixed? what are the consequences?

    randomContent := make([]byte, length)
    _, err = rand.Read(randomContent)
    if err != nil {
        log.Fatal(err)
    }

    // now we will build our result
    result := make([]byte, 2)
    result[0] = 0x02
    result[1] = byte(length)

    // and now we add on the contents
    // randomContent... means that that slice is enumerated since append
    // does not take another slice as the argument
    result = append(result, randomContent...)

    return result
}

func main() {

}


