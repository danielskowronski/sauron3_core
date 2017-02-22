# sauron3
third approach to Sauron - a real time eye on your network - this time in Go Language 

## features
defined in config hosts have ports (tcp/udp/icmp) that are pinged on demand and displayed in webpage - ideal to be put on plasma tv in Networks Operation Command Center or other similar facilites or accessed from mobile phone 

## dev - dependencies
 - `go get gopkg.in/yaml.v2`


## dev - compilation
`go generate && go build -o sauron.o && ./sauron.o`
`sysctl -w net.ipv4.ping_group_range="0 0"` - setup Linux OS for using ping