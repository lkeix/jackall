go build -o ./testdata/application/jackall cmd/jackall/main.go

cd ./testdata/application

./jackall ./...
