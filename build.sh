mkdir -p build
output="palworld-guard"
if [ "$GOOS" == "windows" ]
then
output="$output.exe"
fi
go build -o "./build/$output" -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"