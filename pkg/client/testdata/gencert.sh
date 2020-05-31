cfssl gencert -initca ca_csr.json |cfssljson -bare ca
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config signing.json -profile client csr.json |cfssljson -bare client
cfssl gencert -ca ca.pem -ca-key ca-key.pem -config signing.json -profile server csr.json |cfssljson -bare server
rm *.csr
