cd ../client
GOOS=js GOARCH=wasm go build -o ../wasm/main.wasm
cd ../wasm

#GOOS=js GOARCH=wasm go build -o .
 #mv wasm ./webroot/main.wasm