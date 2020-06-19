<p align="center">
    <img src="img/logo.png" alt="signer" title="signer" class="img-responsive" />
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/lulucas/go-signer?tab=doc"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="godoc" title="godoc"/></a>
    <a href="https://github.com/lulucas/go-signer/releases"><img src="https://img.shields.io/github/v/tag/lulucas/signer" alt="semver tag" title="semver tag"/></a>
    <a href="https://goreportcard.com/report/github.com/lulucas/go-signer"><img src="https://goreportcard.com/badge/github.com/lulucas/go-signer" alt="go report card" title="go report card"/></a>
    <a href="https://github.com/lulucas/go-signer/blob/master/LICENSE"><img src="https://img.shields.io/github/license/lulucas/signer" alt="license" title="license"/></a>
</p>

<p align="center">
    <a href="./README.md">Documentation</a> | 
    <a href="./README_zhCN.md">中文文档</a>
</p>

# Singer

Signer is a golang package which generates sign string for api request。

Generally for payment api or auth request.

## Features

- Support hash functions，MD5、SHA-1 or custom hash function
- Can custom join charactor
- Can custom key value join function
- Support struct and read key name from specific tag
- Support map，url.Values

## Default Options

1. Sort keys at asc
1. Skip empty value
1. Using `MD5` hash
1. Skip `sign` key field
1. Default join charactor `&`
1. Default key value join function `${key}=${value}`
1. Default post function `${str}${joinChar}$keyPairFunc(key, ${key})`

## Example

```go
package main

import (
	"fmt"
	"github.com/lulucas/go-signer"
	"github.com/lulucas/go-signer/hash"
	"net/url"
)

func main() {
	// Generally, you just use Sign function, StrToSign is exported for easy debugging.

	// print amount=1&subject=test&key=123
	fmt.Println(signer.New().Key("123").StrToSign(map[string]interface{}{
		"amount":  1,
		"subject": "test",
	}))

	type Request struct {
		Amount  int    `json:"amount"`
		Subject string `json:"subject,omitempty"`
	}
	// print amount=1&subject=&key=123
	fmt.Println(signer.New().Key("123").NoSkipEmpty().Tag("json").StrToSign(Request{
		Amount:  1,
		Subject: "",
	}))

	// print amount+1#subject+test123
	values := url.Values{
		"amount":    []string{"1"},
		"subject":   []string{"test"},
		"empty":     nil,
		"empty2":    []string{},
		"ignore_me": []string{"ignore"},
	}
	fmt.Println(signer.New().
		Key("123").
		IgnoreKeys("ignore_me").
		JoinChar("#").
		PostHookFunc(func(s, joinChar, key string, kvJoinFn signer.KvJoinFn) string {
			return s + key
		}).
		KvJoinFunc(func(key, value string) string {
			return fmt.Sprintf("%s+%s", key, value)
		}).StrToSign(values))

	// print 53e13a80fedc59e319fdd632caa1c243
	fmt.Println(signer.New().Key("123").Tag("json").Sign(Request{
		Amount:  1,
		Subject: "test",
	}))

	// Use sha1 hash function and make result to upper case.
	// print A00F113DDB21C7B09F305D8855B4FE36E62C0BE1
	fmt.Println(signer.New().Key("123").Tag("json").HashFunc(hash.SHA1(true)).Sign(Request{
		Amount:  1,
		Subject: "test",
	}))
}
```