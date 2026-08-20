package crypto

import "testing"

func TestEnvelopeCodecRoundTrip(t *testing.T) {
	codec, err := NewEnvelopeCodec("controller-c")
	if err != nil {
		t.Fatal(err)
	}
	secret := []byte("0123456789abcdef")
	ciphertext, digest, err := codec.Seal("sensor-c", 1, secret)
	if err != nil {
		t.Fatal(err)
	}
	opened, err := codec.Open("sensor-c", 1, ciphertext)
	if err != nil || string(opened) != string(secret) {
		t.Fatalf("round trip failed: %s %v", opened, err)
	}
	if err := ValidateDigest(opened, digest); err != nil {
		t.Fatal(err)
	}
}
