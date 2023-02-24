go build -o ./replace cmd/replace_singlechecker/replace_singlechecker.go

go build -overlay=`./replace` -o ./testdata/application/jackall cmd/jackall/main.go

cd ./testdata/application

./jackall ./...
# go vet -vettool=./jackall ./...
