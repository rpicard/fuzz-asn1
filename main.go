package main

import (
    "crypto/rand"
    "os"
    "log"
    "path"
    "math/big"
    "fmt"
    "github.com/jawher/mow.cli"
)

type EncodingRuleset interface {
    RandomBoolean() []byte
    RandomInteger() []byte
    RandomBitString() []byte
    RandomOctetString() []byte
    RandomNull() []byte
    RandomObjectIdentifier() []byte
    //RandomReal() []byte
    //RandomEnumerated() []byte
}

func main() {

    app := cli.App("fuzz-asn1", "generate random asn.1 data")

    app.Action = func() {
        cwd, err := os.Getwd()
        if err != nil {
            log.Fatal(err)
        }

        err = os.Mkdir(path.Join(cwd, "tmp"), 0777)
        if err != nil {
            log.Fatal(err)
        }

        // just generating one value for now
        outFile, err := os.Create(path.Join(cwd, "tmp", "object_identifier.asn1"))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", BerEncoding.RandomObjectIdentifier())
    }

    app.Run(os.Args)
}

func GetRandomContent() []byte {

    // we are just going to generate a random length and random content of that
    // length

    // get a random int in [0, 128)
    bigLength, err := rand.Int(rand.Reader, big.NewInt(128))
    if err != nil {
        log.Fatal(err)
    }

    // convert to a regular old int
    length := int(bigLength.Int64())

    randomContent := make([]byte, length)
    _, err = rand.Read(randomContent)
    if err != nil {
        log.Fatal(err)
    }

    return randomContent
}


