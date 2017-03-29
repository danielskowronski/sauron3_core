# sauron3
third approach to Sauron - a real time eye on your network - this time in Go Language 

![demo](https://sc-cdn.scaleengine.net/i/22e315caecf77506e50be0619e57303e.png)

## features
defined in config hosts have ports (tcp/~~udp~~/icmp) that are pinged on demand and displayed in self refreshing webpage - ideal to be put on plasma tv in Networks Operation Command Center or other similar facilites or accessed from mobile phone 

## run
 - config.yml file in same dir as executable (see example)
 - browser -> http://localhost:8888 (port hardcoded)

## dev - dependencies 
 - `go get gopkg.in/yaml.v2`
 - `go get golang.org/x/net/icmp`
 - `go get golang.org/x/net/ipv4`

## dev - compilation
`go generate && go build -o sauron.o && ./sauron.o` 

`sysctl -w net.ipv4.ping_group_range="0 0"` - setup Linux OS for using ping
