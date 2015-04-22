package main

import (
    "crypto/rand"
    "log"
)

type berEncoding struct {}

var BerEncoding berEncoding

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

    result[2] = randomByte[0]

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
    //
    // now we generate our random content
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

    randomContent := GetRandomContent()

    // now we will build our result
    result := make([]byte, 2)
    result[0] = 0x02
    result[1] = byte(len(randomContent))

    // and now we add on the contents
    // randomContent... means that that slice is enumerated since append
    // does not take another slice as the argument
    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomBitString() []byte {

    // the encoding of a bitstring can be either primitive or constructed
    // we are going to stick to primitive for now
    //
    // TODO: generate constructed bitstrings as well
    //
    // the plan is to generate a random length and random content of that
    // length
    //
    // there is some special meaning to the first octet of the content:
    //
    //   8.6.2.2 The initial octet shall encode, as an unsigned binary
    //   integer with bit 1 as the least significant bit, the number of
    //   unused bits in the final subsequent octet. The number shall be in
    //   the range zero to seven. 
    //
    // we do not need to deal with that, as our random bytes will randomly
    // tell the parser how many bits are unused in the last octet
    //
    // the length includes that special first octet 

    randomContent := GetRandomContent()

    // now we will build our result
    result := make([]byte, 2)
    result[0] = 0x03
    result[1] = byte(len(randomContent))

    // and now we add on the contents
    // randomContent... means that that slice is enumerated since append
    // does not take another slice as the argument
    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomOctetString() []byte {

    // the encoding of an octetstring can be either primitive or constructed
    // we are going to stick to primitive for now
    //
    // TODO: generated constructed octetstrings as well

    randomContent := GetRandomContent()

    // now we will build our result
    result := make([]byte, 2)
    result[0] = 0x04
    result[1] = byte(len(randomContent))

    // and now we add on the contents
    // randomContent... means that that slice is enumerated since append
    // does not take another slice as the argument
    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomObjectIdentifier() []byte {

    // the object identifier seems to just be a basic tag, length, value
    // encoding
    //
    // it looks like we might not be generating a totally complete
    // identifier
    //
    // TODO: make sure the identifiers that are generated are valid

    randomContent := GetRandomContent()

    result := make([]byte, 2)
    result[0] = 0x06
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomNull() []byte {

    // there is nothing to randomize
    // this is always the same value
    return []byte{0x05, 0x00}
}

func (berEncoding) RandomReal() []byte {

    // there are some special values that might be interesting
    //
    // one octet:
    //
    //    0100000 - PLUS-INFINITY
    //    0100001 - MINUS-INFINITY
    //
    // in theory they can be generated at random by this code, but
    // it might be worth generating those specifically
    //
    // TODO: look into that ^

    randomContent := GetRandomContent()

    result := make([]byte, 2)
    result[0] = 0x09
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomEnumerated() []byte {

    // this is supposedly encoded in EXACTLY the same way as the integer
    // this means that the same length byte weirdness applies here too
    // if we implement it for integer we should implement it here too

    randomContent := GetRandomContent()

    result := make([]byte, 2)
    result[0] = 0x0A
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomNumericString() []byte {

    charSet := "0123456789 "

    randomContent := GetRandomContentFromCharset(charSet)

    result := make([]byte, 2)
    result[0] = 0x12
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomPrintableString() []byte {

    // i do not like having this on one long line
    charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890 '()+,-./:=?"

    randomContent := GetRandomContentFromCharset(charSet)

    result := make([]byte, 2)
    result[0] = 0x13
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}

func (berEncoding) RandomIa5String() []byte {

    // this is an ascii string
    // we are just going to generate a character set from all bytes < 128
    var charSet string

    for i := 0; i < 128; i++ {
        charSet += string(i)
    }

    randomContent := GetRandomContentFromCharset(charSet)

    result := make([]byte, 2)
    result[0] = 0x16
    result[1] = byte(len(randomContent))

    result = append(result, randomContent...)

    return result
}
