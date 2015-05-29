# go-codeivate
Eu não tinha nada pra fazer, não queria ir pro Facebook, então fiz isso.


Clone my repo

```
git clone https://github.com/schleumer/go-codeivate.git
```

`cd` on it

```
cd go-codeivate
```

Get dependencies

```
go get github.com/nsf/termbox-go
go get github.com/gizak/termui
```

# Running

### On cool OSes

``` 
go run main.go -username your_nice_username
```

### If you are using Windows you need to execute that way

```
go build main.go
start "" "main.exe" -username your_nice_username
```

This will avoid termbox to overwrite `cmd`'s settings(yeah)


# NOTICE ME SENPAI