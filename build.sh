output="palworld-guard"
if [ "$GOOS" == "windows" ]
then
output="$output.exe"
fi
go build -o output -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"