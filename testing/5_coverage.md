1. command to cover
cd eval
go test -run=Coverage -coverprofile=c.out
go tool cover -html=c.out

2. code that are not covered (not read during test) are in read.