## Installation

```bash
go get github.com/NovanHsiu/goutil
```

## Usage

```go
package main

import (
	"github.com/NovanHsiu/goutil"
	"fmt"
	"time"
)

func main() {
	fmt.Println(goutil.SQLTimeStringToTime(time.Now()))
}
```