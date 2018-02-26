# Geetest SDK for Golang

*This SDK is made by personal, not the official version*

# Installation
```
go get -u github.com/Secbone/geetest
```

# Usage

```golang
import "github.com/Secbone/geetest"

tester := geetest.New(ID, KEY)

info := tester.Register()

result := tester.Validate(fallback, challenge, validate, seccode)
```
