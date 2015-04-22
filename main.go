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

        GenerateLotsOfSamples(encoding, path.Join(cwd, "tmp"))
    }

    app.Run(os.Args)
}

func GenerateLotsOfSamples(er EncodingRuleset, outdir string) {

    // could have the user specify a number of samples on the command line
    // not sure if that really matters. for now we are just going to generate
    // a relatively appropriate amount of each type
    //
    // having a loop for each type feels pretty verbose. there might be a
    // better way to do this, but for now it gets the job done

    // generate 10 booleans
    for i := 0; i < 10; i++ {
        fileName := fmt.Sprintf("rand_bool_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomBoolean())
        outFile.Close()
    }

    // generate 1000 integers
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_integer_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomInteger())
        outFile.Close()
    }

    // generate 1000 bitstrings
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_bitstring_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomBitString())
        outFile.Close()
    }

    // generate 1000 octetstrings
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_octetstring_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomOctetString())
        outFile.Close()
    }

    // generate 1 null
    // still looping for consistency
    for i := 0; i < 1; i++ {
        fileName := fmt.Sprintf("rand_null_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomNull())
        outFile.Close()
    }

    // generate 1000 object identifiers
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_object_identifier_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomObjectIdentifier())
        outFile.Close()
    }

    // generate 1000 reals
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_real_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomReal())
        outFile.Close()
    }

    // generate 1000 enumerateds
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_enumerated_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomEnumerated())
        outFile.Close()
    }

    // generate 1000 numeric strings
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_numeric_string_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomNumericString())
        outFile.Close()
    }

    // generate 1000 printable strings
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_printable_string_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomPrintableString())
        outFile.Close()
    }

    // generate 1000 ia5 strings
    for i := 0; i < 1000; i++ {
        fileName := fmt.Sprintf("rand_ia5_string_%d.asn1", i)

        outFile, err := os.Create(path.Join(outdir, fileName))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Fprintf(outFile, "%s", er.RandomIa5String())
        outFile.Close()
    }
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
