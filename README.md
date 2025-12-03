# IP Resolver

IP Resolver is a small Golang library that resolves an IP to an ASN and a country. It is useful for determining the language a user's system may be using.

## Usage

You'll first need an IP to ASN database. The format must be:

| Key | Type | Example |
| --- | ---- | ------- |
| `range_start` | string | 1.0.0.0 |
| `range_end` | string | 1.0.0.255 |
| `AS_number` | int | 13335 |
| `country_code` | string | US |
| `AS_description` | string | CLOUDFLARENET |

1. Load the database:

```go
// load the asn database
err := ipresolve.LoadIPFile("./test.tsv")
if err != nil {
    panic(err)
}
```

2. Query the database:

```go
// resolve the IP
data, err := ipresolve.ResolveIp(ip)
if err != nil {
    panic(err)
}

// returns &{1.1.1.0 1.1.1.255 US 13335 CLOUDFLARENET}
```

## Database

The database can be downloaded from [iptoasn.com](https://iptoasn.com). Thanks to Frank Denis for providing it.