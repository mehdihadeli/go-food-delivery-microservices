# Certs

## CA Certs

```bash
mkdir -p deployments/kustomize/certs
# Generate CA key
# Use pass-phrase: test
openssl genrsa -des3 -out deployments/kustomize/certs/ca.key 4096
# Generate Root CA crt from key
openssl req -x509 -new -days 1825 -key deployments/kustomize/certs/ca.key -out deployments/kustomize/certs/ca.crt
```

> Prompts

```
Country Name (2 letter code) []:US
State or Province Name (full name) []:CA
Locality Name (eg, city) []:Riverside
Organization Name (eg, company) []:Sumo
Organizational Unit Name (eg, section) []:Demo
Common Name (eg, fully qualified host name) []:*
Email Address []:ca@sumo.com
```

## Micro Certs

```bash
# Create microservice key
openssl genrsa -out deployments/kustomize/certs/micro.key 2048
# Generate CSR
openssl  req -new -key deployments/kustomize/certs/micro.key -out deployments/kustomize/certs/micro.csr
```

> Prompts, Provide `empty` challenge password

```
Country Name (2 letter code) []:US
State or Province Name (full name) []:CA
Locality Name (eg, city) []:Riverside
Organization Name (eg, company) []:Sumo
Organizational Unit Name (eg, section) []:Demo
Common Name (eg, fully qualified host name) []:*
Email Address []:micro@sumo.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
```

### Sign CSR with CA

```bash
openssl x509 -req -days 365 -in deployments/kustomize/certs/micro.csr -CA deployments/kustomize/certs/ca.crt -CAkey deployments/kustomize/certs/ca.key -CAcreateserial -out deployments/kustomize/certs/micro.crt
```
