# tinygo-flash

Still in development.  

This is a repository for fixing the problem of tinygo flash not working on Japanese windows.  
Until it's merged into tinygo, it will serve as a workaround for the private version of the above problem.  

So far, I've only tested it in the following environments.  

* adafruit PyPortal
* adafruit feather-m4

## Usage

```
$ tinygo-flash --port COM7 --target feather-m4 your_application.uf2
```

## Installation

```
go get github.com/sago35/tinygo-flash
```

## Build

```
$ go build
```

or

```
$ go run dist/make_dist.go 1.2.3
```

or

```
$ go run dist/make_dist.go 1.2.3 --dev
```

### Environment

* go
* kingpin.v2

## Notice

## Author

sago35 - <sago35@gmail.com>
