package certs

import _ "embed"

//go:embed "ca-cert.pem"
var CACertPEMBlock []byte

//go:embed "server-key.pem"
var ServerKeyPEMBlock []byte

//go:embed "server-cert.pem"
var ServerCertPEMBlock []byte
