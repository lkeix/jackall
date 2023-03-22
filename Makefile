build:
	go install github.com/lkeix/jackall/cmd/replace_singlechecker@latest
	go build -overlay=`./replace` -o ./jackall cmd/jackall/main.go
	mv ./jackall $$GOPATH/bin
