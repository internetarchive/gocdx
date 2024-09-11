# gocdx

`gocdx` is a Go package for parsing CDX files. CDX files are commonly used in web archiving to index the contents of WARC (Web ARChive) files. This package is maintained by the Internet Archive.

## Installation

To use this package in your Go project, you can install it using `go get`:

```bash
go get github.com/internetarchive/cdx
```

## Usage (as a package)

Here's a basic example of how to use the cdx package:

```go
package main

import (
	"fmt"
	"os"

	"github.com/internetarchive/cdx"
)

func main() {
	file, err := os.Open("path/to/your/cdx/file.cdx")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	records, err := cdx.Parse(file)
	if err != nil {
		fmt.Println("Error parsing CDX file:", err)
		return
	}

	fmt.Printf("Parsed %d records\n", len(records))
	for i, record := range records {
		fmt.Printf("Record %d: %+v\n", i+1, record)
	}
}
```

## Original Author

This project was originally created by [Corentin Barreau](https://github.com/CorentinB) at the Internet Archive.

## License
This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).
See the LICENSE file for details.

## Contact
If you have any questions or feedback, please open an issue on the GitHub repository at https://github.com/internetarchive/cdx.
For more information about the Internet Archive, visit https://archive.org/.
