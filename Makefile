

name := aflex-du
version := 0.0.1


build:

	env GOOS=darwin GOARCH=amd64  go build -o ./bin/macos-amd64/$(name) .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/macos-arm64/$(name) .
	env GOOS=windows GOARCH=386  go build -o ./bin/windows-386/$(name) .
	env GOOS=windows GOARCH=amd64  go build -o ./bin/windows-amd64/$(name).exe .
	env GOOS=linux GOARCH=amd64  go build -o ./bin/linux-amd64/$(name) .
	env GOOS=linux GOARCH=arm64  go build -o ./bin/linux-arm64/$(name) .
	env GOOS=linux GOARCH=arm  go build -o ./bin/linux-arm/$(name) .

	# Create the README.md file automatically
	./build-readme.sh > README.md

	# just for local machine convenience
	cp ./bin/linux-amd64/$(name) ../Flex_Prestel/tools

clean:
	rm -R  ./bin/*



