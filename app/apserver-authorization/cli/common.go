package cli

var (
	privateKeyBytes = []byte(`-----BEGIN RSA Private Key-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMpaaRtyJDvmVKFs
0EdaP3gwVvKesY5Qj1iTs3CKsh5kCp0HKhQhmZW5zV3X727p9JY9EeBcXMYjU7ET
HGL9kY+AcV52NsT6/lrWbrop0lJX2w7MXZghCxd/E50cY89XGR0GvNVG8ee212qr
duohURrKaWbyavo+uiNd2aN4EUgTAgMBAAECgYEAoNDBf6Jy0Xf4AmJsFIjQsEAa
ma8tBSFZCtg3X1WawTRYivtKob0iRi/n2pDmJIuialQWhOxifsVMmgnKIZHLv4o/
2SeMkklilg9SIhIIiPD+xubbE3eVQGUcZPmqAxSk/E4NaXNC7coSGZCiUrGeCJfR
3J2HQEMzDarWCFnmp6ECQQDeDBPOZ3oXeAaQ8ZyXIxwwi2ntKh3YTb5R55GH7T12
UO1TNK2LzTfgIUXSn/rH5/poq0RyxLBSaqcNzV+OhMJvAkEA6Utr2HEC3NDxlfVa
Q6fnIVr4ETo6v3JN3Cqtq8kZ5XX7HTLo1WhulJf3fxHwh3jD8s5LKN5aEm+HLPmX
fVeWnQJAFPV64R6vTYvMwt2rdDCiNorSQsqY6pPcBQsgl33zMTnOTO5J+0oxnfxG
BO2I1Fm3Ly4LVfHu2riqcAkUnfU2DQJAAojXOxq/NTbv6PkpaeLBGBOs7kL7sGjF
f8bW7C7bISsO91o+PVNNIEAmaDMBsfcV6eVj26XOxLSBe3Oaubnh4QJBANNdRH9B
OngI5bfNhUKNC9zReonOecZovJ6tAoPfA1oxoHrUFSt5wdUxqwaMUp9pydLODPqN
eUgUzfSmGnGD0/4=
-----END RSA Private Key-----`)

	publicKeyBytes = []byte(`-----BEGIN RSA Public Key-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDKWmkbciQ75lShbNBHWj94MFby
nrGOUI9Yk7NwirIeZAqdByoUIZmVuc1d1+9u6fSWPRHgXFzGI1OxExxi/ZGPgHFe
djbE+v5a1m66KdJSV9sOzF2YIQsXfxOdHGPPVxkdBrzVRvHnttdqq3bqIVEaymlm
8mr6ProjXdmjeBFIEwIDAQAB
-----END RSA Public Key-----
`)
)

func getPrivateKey() []byte {
	return privateKeyBytes
}

func getPublicKey() []byte {
	return publicKeyBytes
}
