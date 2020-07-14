1. if GOPATH is not set, we must export PATH for gopherjs binary after go get command. in this case is /home/<user-name>/go/bin. add it in .bashrc

2. go version is incompatible with go1.14. so we must let gopherjs use go1.12.6:
```sh
go get golang.org/dl/go1.12.16
go1.12.16 download
export GOPHERJS_GOROOT="$(go1.12.16 env GOROOT)" # Also add this line to your .profile or equivalent.
```

3. node version 8.11.3 can be used