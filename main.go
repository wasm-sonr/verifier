package main

import (
	"crypto/ecdsa"
	"crypto/sha3"

	"github.com/extism/go-pdk"
	"github.com/sonr-io/crypto/mpc"
)

type VerifyRequest struct {
	PubKey  []byte `json:"pub_key"`
	Message []byte `json:"message"`
	Sig     []byte `json:"sig"`
}

type VerifyResponse struct {
	Valid bool `json:"valid"`
}

func main() {
	verify()
}

//go:wasmexport verify
func verify() int32 {
	req := VerifyRequest{}
	err := pdk.InputJSON(req)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1
	}
	pdk.Log(pdk.LogInfo, "Deserialized request successfully")
	res := VerifyResponse{Valid: false}
	valid, err := VerifyWithPubKey(req.PubKey, req.Message, req.Sig)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	res.Valid = valid
	pdk.OutputJSON(res)
	return 0
}

func VerifyWithPubKey(pubKeyCompressed []byte, data []byte, sig []byte) (bool, error) {
	edSig, err := mpc.DeserializeSignature(sig)
	if err != nil {
		return false, err
	}
	ePub, err := mpc.GetECDSAPoint(pubKeyCompressed)
	if err != nil {
		return false, err
	}
	pk := &ecdsa.PublicKey{
		Curve: ePub.Curve,
		X:     ePub.X,
		Y:     ePub.Y,
	}

	// Hash the message using SHA3-256
	hash := sha3.New256()
	hash.Write(data)
	digest := hash.Sum(nil)
	return ecdsa.Verify(pk, digest, edSig.R, edSig.S), nil
}
