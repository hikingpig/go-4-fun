1. net/url uses net/http for testing while net/http import net/url. this creates an import cycle.
- we must declare a file called url_test.go. that file can import net/http and net/url to run the test required. it is not included in `go build` command

* this example is not correct. url_test.go is not an external test file!

- the external test files at the same level of package's files. but they have package name with suffix _test, like fmt_test


2. GoFiles list production code, included in `go build`
go list -f={{.GoFiles}} fmt
go list -f={{.GoFiles}} net/http
go list -f={{.GoFiles}} net/url

TestGoFiles includes package's test files
go list -f={{.TestGoFiles}} fmt
go list -f={{.TestGoFiles}} net/url

3. the isSpace func defined it fmt package, to avoid using unicode.IsSpace method. However, we want this isSpace works the same as unicode.IsSpace so we do a test for isSpace.
- but fmt_test, because it is not the same package with fmt, can not access isSpace. so we declare a file call export_test.go, that defines the variable IsSpace to be accessible only for testing.
So:
  - the unicode is only imported in testing, not in build
  - the isSpace is only imported in testing, not in build

