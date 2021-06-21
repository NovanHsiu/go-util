## Installation

```bash
GOPATH=`go env | grep GOPATH | cut -d = -f 2 | sed 's/"//g'`
cd $GOPATH/src
mkdir acebiotek
cd acebiotek
git clone http://10.11.101.98:10080/huey_yu/simplecipher.git
```

## Usage

```go
package main

import (
	"acebiotek/simplecipher"
	"fmt"
)

func main() {
	ptext := "abcd"
	pp := "this is a private key"
	entext, _ := simplecipher.Encrypt(ptext, pp)
	fmt.Println(entext)
	decodetext, _ := simplecipher.Decrypt(entext, pp)
	fmt.Println(decodetext)
}
```