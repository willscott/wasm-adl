package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: gen github.com/myorg/mycooladl myadl.wasm\n")
		os.Exit(1)
	}

	pkg := os.Args[1]
	out := os.Args[2]

	// make a temp dir
	dir, err := ioutil.TempDir("", "build")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	// copy the package to the tempdir
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to read dir: %v\n", err)
		os.Exit(1)
	}

	for _, f := range entries {
		if !f.IsDir() {
			if err := copy(path.Join(".", f.Name()), path.Join(dir, f.Name())); err != nil {
				fmt.Printf("Failed to copy %s: %v\n", f.Name(), err)
				os.Exit(1)
			}
		}
	}

	f, err := os.OpenFile(path.Join(dir, "link.go"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to write: %v\n", err)
		os.Exit(1)
	}

	templ := fmt.Sprintf(`
package main

import (
	theADL "%s"
)

func init() {
	ADL = theADL.Reify
}
`, pkg)

	if _, err := f.WriteString(templ); err != nil {
		fmt.Printf("Failed to write: %v\n", err)
		os.Exit(1)
	}
	f.Close()

	// get go mod set up
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to setup module: %v\n", err)
		os.Exit(1)
	}

	// get the wasm
	cmd = exec.Command("go", "generate", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build: %v\n", err)
		os.Exit(1)
	}

	// stream back to output
	contents, err := os.ReadFile(path.Join(dir, "adl.wasm"))
	if err != nil {
		fmt.Printf("Failed to read output: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(out, contents, 0644); err != nil {
		fmt.Printf("Failed to write output: %v\n", err)
		os.Exit(1)
	}
	return
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
