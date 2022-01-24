# databox

`databox` is a tool converting file or directory to byte buffer and storing in go file.

## install

`go install github.com/snowmerak/databox@latest`

## how to use

`intogo <target-directory-or-file> <package-name> <variable-name>`

1. read all target-directory-or-file and that's sub directorys and files to buffer
2. create ./package-name/variable-name.go
3. write buffer to ./package-name/variable-name.go
