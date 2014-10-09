dnsmadeeasy
===========

DNS Made Easy Golang Api

This API is a work in progress and is not production ready.

Please feel free to report bugs, fix bugs, implement missing API calls and do pull requests.

I'm also new to go, so please feel free to comment on the code itself, pointing out flaws, non idiomatic go, and other mistakes!

# How to install?

`go get github.com/huguesalary/dnsmadeeasy`

# How to use?

```go
package main

import (
	"github.com/huguesalary/dnsmadeeasy"
)

func main() {
	client := dnsmadeeasy.NewClient("yourapikey", "yourapisecret")
	domain := client.CreateDomain("yourdomain.com")

	record := &dnsmadeeasy.Record{}
	record.Type = "A"
	record.Name = "a-record"
	record.Value = "127.0.0.1"
	record.GtdLocation = "DEFAULT"
	record.Ttl = 86600
	client.AddRecord(domain.Id, rec)
}
```
# Known Issue(s)

- Because the DNS Made Easy test API SSL certificates are not valid, the code forcefully disable SSL verification. This will be fixed soon.