# my-2fa
[![Go Report Card](https://goreportcard.com/badge/github.com/noxtal/my-2fa)](https://goreportcard.com/report/github.com/noxtal/my-2fa)

A simple 2FA Backend API entirely in Go as a way to learn the language. Cryptographically secure using `crypto/rand` according to a [gosec](https://github.com/securego/gosec) scan. Manages a database internally to handle multiple users and codes. This should still not be used in production as it is just a quick learning project and was not built professional.

## Usage
### GET
Get with query `user` to get a new random code for the user specified.

- `200`: A code was generated and sent to the user.
- `500`: The user parameter was missing.

### POST
Post with query `user` and `code` to validate a code for the user specified.

- `200`: Access Granted.
- `401`: Access Denied.
- `500`: One of the parameters was missing.

## TODO
- [x] Lay out the backend
- [x] Use a password to validate identity
- [ ] Wiping every minute
- [ ] Frontend