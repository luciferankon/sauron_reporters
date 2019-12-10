module github.com/step/sauron_reporters

go 1.12

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/step/angmar v0.0.0-20191127113211-fbeaab94f9b7
	github.com/step/saurontypes v0.0.0-20191127114135-1c7b69a4e64f
	github.com/step/uruk v0.0.0-20191127114036-eb84283fad8d
	github.com/tidwall/pretty v1.0.0 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.1.3
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/step/uruk => ../../step/uruk/

replace github.com/step/saurontypes => ../../step/saurontypes/

replace github.com/step/angmar => ../../step/angmar/
