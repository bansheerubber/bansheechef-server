package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	directory, _ := os.Getwd()
	directory = directory + "/" + os.Args[1] + "_type.gen.go"

	config := &packages.Config{
		Mode: packages.LoadSyntax,
		Tests: false,
	}
	packages, _ := packages.Load(config)

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "package %s\n\n", packages[0].Name)
	fmt.Fprint(&buffer, "import \"reflect\"\n\n")
	fmt.Fprintf(&buffer, "func %s_type() reflect.Type {\n", os.Args[1])
	fmt.Fprintf(&buffer, "\treturn reflect.TypeOf((*%s)(nil)).Elem()\n", os.Args[1])
	fmt.Fprint(&buffer, "}\n\n")

	ioutil.WriteFile(directory, buffer.Bytes(), 0777)
}
