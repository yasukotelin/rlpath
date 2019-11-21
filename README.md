# rlpath.go

rlpath is scanning path with completion library.

![rlpath demo](./images/rlpath-demo.gif)

## Instalation

```
$ go get -u github.com/yasukotelin/rlpath
```

```go
import (
	"github.com/yasukotelin/rlpath"
)
```

## Example

You can use this library so easily!!

```
package main

import (
	"fmt"
	"log"

	"github.com/yasukotelin/rlpath"
)

func main() {
	scanner := rlpath.Scanner{
		Prompt:  "$ ",
		RootDir: "~/go",
		OnlyDir: false,
	}

	path, err := scanner.Scan()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(path)
}
```

`rlpath.Scanner` has options.

| Option  |                                                                                                    |
|---------|----------------------------------------------------------------------------------------------------|
| Prompt  | Prompt is left edge text, like a $.                                                                |
| RootDir | RootDir is root directory to start scanning. If this is empty, start scanning from execution path. |
| OnlyDir | OnlyDir is flag that shows only directory.                                                         |

## Author

yasukotelin

## LICENCE

MIT
