package models

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

func SplitToken(token string) []string {
	return strings.Split(token, ".")
}

func ParseHeader(headerString string) (*Header, error) {
	var header Header
	decodedStringFromBase64, err := base64.RawURLEncoding.DecodeString(headerString)
	if err != nil {
		return nil, err
	}
	err = ParseJSON(&header, string(decodedStringFromBase64))
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func ParsePayload(payloadString string) (*Payload, error) {
	var payload Payload
	decodedStringFromBase64, err := base64.RawURLEncoding.DecodeString(payloadString)
	if err != nil {
		return nil, err
	}
	err = ParseJSON(&payload, string(decodedStringFromBase64))
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func ParseJSON(target JSONSerializable, jsonString string) error {
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	err := decoder.Decode(target)
	if err != nil {
		return err
	}
	return nil
}
