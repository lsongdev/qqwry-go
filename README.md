# qqwry

> 🌎 A lightweight native Golang implementation of GeoIP API


## Example

```go
package cli

import (
	"fmt"

	"github.com/song940/qqwry/qqwry"
)

func main() {
	q, err := qqwry.NewQQwry("qqwry.dat")
	if err != nil {
		panic(err)
	}
	result, err := q.Find("1.1.1.1")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.IP)
	fmt.Println(result.Country, result.City)
}

```

## License

This project is licensed under the MIT license.