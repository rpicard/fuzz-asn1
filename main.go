package main

import (
    "crypto/rand"
    "errors"
    "strings"
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
    RandomReal() []byte
    RandomEnumerated() []byte
    RandomNumericString() []byte
    RandomPrintableString() []byte
    RandomIa5String() []byte
}

func main() {

    app := cli.App("fuzz-asn1", "generate random asn.1 data")

    encodingOpt := app.StringOpt("e encoding", "ber", "which encoding ruleset should we use?")

    app.Action = func() {

        var encoding EncodingRuleset

        switch strings.ToLower(*encodingOpt) {
        case "ber":
            encoding = BerEncoding
        default:
            log.Fatal(errors.New("I do not know that encoding."))
        }

        cwd, err := os.Getwd()
        if err != nil {
            log.Fatal(err)
        }

        err = os.Mkdir(path.Join(cwd, "tmp"), 0777)
        if err != nil {
            log.Fatal(err)
        }

        // just generating one value for now
        outFile, err := os.Create(path.Join(cwd, "tmp", "output.asn1"))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", encoding.RandomIa5String())
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

func GetRandomContentFromCharset(charSet string) []byte {

    // get a random length
    bigLength, err := rand.Int(rand.Reader, big.NewInt(128))
    if err != nil {
        log.Fatal(err)
    }

    // convert to a regular old int
    length := int(bigLength.Int64())

    result := make([]byte, length)

    for i := 0; i < length; i++ {
        // get a random character from the charset
        bigIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
        if err != nil {
            log.Fatal(err)
        }

        index := int(bigIndex.Int64())
        result[i] = charSet[index]
    }

    return result
}
