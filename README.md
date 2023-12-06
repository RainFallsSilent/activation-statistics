## activation statistics of ELA and ESC

## go version >= 1.18

## install
```shell
go get git@github.com:RainFallsSilent/activation-statistics.git
```

## make
```shell
go mod tidy
make all
```

## run
```shell
./activation
```

```shell
./activation /home/config.json
```

default config.json:
```json
{
	"Days":      2,
	"StartHour": 8,
	"ELARpcConfig": {
		"HttpUrl": "https://api.elastos.io/ela",
		"User":    "",
		"Pass":    ""
	},
	"ESCRpcConfig": {
		"HttpUrl": "https://api.elastos.io/esc",
		"User":    "",
		"Pass":    ""
	}
}
```