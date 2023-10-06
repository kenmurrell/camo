# CAMO

An image steganography encoder written in Golang.

### USE

To encode a file into a picture:
```Bash
./camo encode -hide {file-to-encode} -host {host-file} [-blue] [-encrypt]
```

Similarly, to decode a file from a host picture:
```Bash
./camo decode -host {host-file} -output {output-filename} [-blue] [-decrypt]
```
