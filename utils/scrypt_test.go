package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestHashingPlaintext(t *testing.T) {
    s, err := NewScrypt()
    assert.NotNil(t, s)
    assert.Nil(t, err)

    hash, err := s.HashPlaintext("my lil plaintext")
    assert.NotEmpty(t, hash)
    assert.Contains(t, hash, "32768@8@1")
    assert.Nil(t, err)
}

func TestMatchingPlaintext(t *testing.T) {
    s, err := NewScrypt()
    assert.NotNil(t, s)
    assert.Nil(t, err)

    hash, err := s.HashPlaintext("my lil plaintext")
    assert.NotEmpty(t, hash)
    assert.Contains(t, hash, "32768@8@1")
    assert.Nil(t, err)

    s2, err := LoadScryptFromHash(hash)
    assert.NotNil(t, s2)
    assert.Nil(t, err)

    match, err := s2.MatchesPlaintext("my lil plaintext")
    assert.True(t, match)
    assert.Nil(t, err)
}

func TestBarfOnLoadingGarbage(t *testing.T) {
    s, err := LoadScryptFromHash("asa#$@$%@241e2q 1e!2@22@212")
    assert.Nil(t, s)
    assert.NotNil(t, err)
}

func BenchmarkPlaintextEncryption(b *testing.B) {
    s, _ := NewScrypt()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        s.HashPlaintext("my lil plaintext")
    }
}