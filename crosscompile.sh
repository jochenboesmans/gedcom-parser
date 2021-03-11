declare -a os=("darwin" "darwin" "linux" "windows")
declare -a arch=("amd64" "arm64" "amd64" "amd64")

for i in "${os[@]}"; do
  env GOOS="${os[$i]}" GOARCH="${arch[$i]}" go build -o gedcom-parser."${os[$i]}"."${arch[$i]}"
  7z a gedcom-parser."${os[$i]}"."${arch[$i]}".7z gedcom-parser."${os[$i]}"."${arch[$i]}"
done
