package utils

import "fmt"
import "strconv"
import "strings"
import "crypto/rand"
import "encoding/base64"
import "code.google.com/p/go.crypto/scrypt"

type scryptWrapper struct {
    currentHash string
    salt []byte
    paramN int
    paramR int
    paramP int
}

func NewScrypt() (*scryptWrapper, error) {
    s := &scryptWrapper{}
    err := s.Reset()
    if err != nil {
        return nil, err
    }

    return s, nil
}

func LoadScryptFromHash(hash string) (*scryptWrapper, error) {
    s := &scryptWrapper{}
    err := s.parseFromHash(hash)
    if err != nil {
        return nil, err
    }

    return s, nil
}

func (s *scryptWrapper) parseFromHash(hash string) error {
    // Split the hash and make sure we have all of the components.
    hashPieces := strings.Split(hash, "@")
    if len(hashPieces) != 5 {
        return fmt.Errorf("supplied hash corrupted / has invalid format: expected 5 pieces, got %d", len(hashPieces))
    }

    // Pull out our pieces and parse them.
    paramN, err := strconv.Atoi(hashPieces[0])
    if err != nil {
        return fmt.Errorf("couldn't parse N param; supplied value: %s", hashPieces[0])
    }

    paramR, err := strconv.Atoi(hashPieces[1])
    if err != nil {
        return fmt.Errorf("couldn't parse R param; supplied value: %s", hashPieces[1])
    }

    paramP, err := strconv.Atoi(hashPieces[2])
    if err != nil {
        return fmt.Errorf("couldn't parse P param; supplied value: %s", hashPieces[2])
    }

    salt, err := base64.StdEncoding.DecodeString(hashPieces[3])
    if err != nil {
        return fmt.Errorf("couldn't decode salt; supplied value: %s", hashPieces[3])
    }

    s.paramN = paramN
    s.paramR = paramR
    s.paramP = paramP
    s.salt = salt
    s.currentHash = hash

    return nil
}

func (s *scryptWrapper) Reset() error {
    // Get a new salt here, make sure our parameters are accurate, etc.
    s.paramN = 16384
    s.paramR = 8
    s.paramP = 1
    salt, err := getRandomBytes(8)
    if err != nil {
        return err
    }

    s.salt = salt

    return nil
}

func (s *scryptWrapper) MatchesPlaintext(plaintext string) (bool, error) {
    hashedPlaintext, err := s.HashPlaintext(plaintext)
    if err != nil {
        return false, err
    }

    return constantTimeStringEq(s.currentHash, hashedPlaintext), nil
}

func (s *scryptWrapper) HashPlaintext(plaintext string) (string, error) {
    plaintextRaw := []byte(plaintext)

    encryptedPlaintext, err := scrypt.Key(plaintextRaw, s.salt, s.paramN, s.paramR, s.paramP, 32)
    if err != nil {
        return "", err
    }

    // Now build our param-prepended version of the hash.
    encodedSalt := base64.StdEncoding.EncodeToString(s.salt)
    encodedHash := base64.StdEncoding.EncodeToString(encryptedPlaintext)

    return fmt.Sprintf("%d@%d@%d@%s@%s", s.paramN, s.paramR, s.paramP, encodedSalt, encodedHash), nil
}

func getRandomBytes(n int) ([]byte, error) {
    buf := make([]byte, n)
    _, err := rand.Read(buf)
    if err != nil {
        return nil, err
    }

    return buf, nil
}

func constantTimeStringEq(expected, actual string) bool {
    expectedRaw := []byte(expected)
    actualRaw := []byte(actual)
    expectedLen := len(expectedRaw)
    actualLen := len(actualRaw)
    minLen := actualLen
    if actualLen > expectedLen {
        minLen = expectedLen
    }

    // Compare both strings, doing an XOR.  If the bytes at the given location are equal,
    // they XOR to 0, which is ORd to `result`.  ORing doesn't destroy data, so we can keep
    // ORing and if there was a single byte mismatch, `result` will be != 0 when we're done
    // iterating through the buffers, even if all but one byte matched.
    result := 0
    for i := 0; i < minLen; i++ {
        result |= int(expectedRaw[i]) ^ int(actualRaw[i])
    }

    // Factor in the length of the strings, too, otherwise one could be a prefix of the
    // other and still match.  We do it after the full comparison so we don't break our
    // constant time comparison. e.g. 'abcd' and 'abcdef'
    result |= expectedLen ^ actualLen

    return result == 0
}
