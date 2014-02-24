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
	s, err := LoadScryptFromHash("123")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@456")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@456@789")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@456@789@012")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("asd@lol@wtf@bbq@kfc")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@asd@lol@bbq@kfc")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@456@wtf@bbq@kfc")
	assert.Nil(t, s)
	assert.NotNil(t, err)

	s, err = LoadScryptFromHash("123@456@789@bbq@kfc")
	assert.Nil(t, s)
	assert.NotNil(t, err)
}

func TestErrorsOnUnreasonableParameters(t *testing.T) {
	s, err := NewScrypt()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	// Set a bogus parameter.
	s.paramN = 3312

	e, err := s.HashPlaintext("test plaintext")
	assert.Equal(t, "", e, "ciphertext must be empty")
	assert.NotNil(t, err)
}

func TestErrorsOnUnreasonableParametersFromExistingHash(t *testing.T) {
	s, err := LoadScryptFromHash("1232@1@1@YXNkc2FkYXM=@notarealhash")
	assert.NotNil(t, s)
	assert.Nil(t, err)

	e, err := s.HashPlaintext("test plaintext")
	assert.Equal(t, "", e, "ciphertext must be empty")
	assert.NotNil(t, err)
}

func BenchmarkPlaintextEncryption(b *testing.B) {
	s, _ := NewScrypt()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.HashPlaintext("my lil plaintext")
	}
}
