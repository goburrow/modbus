System testing for [modbus library](https://github.com/goburrow/modbus)

Modbus simulator
----------------
*   [Diagslave](http://www.modbusdriver.com/diagslave.html)
*   [socat](http://www.dest-unreach.org/socat/)

```bash
# TCP
$ diagslave -m tcp -p 5020


# RTU/ASCII
$ socat -d -d pty,raw,echo=0 pty,raw,echo=0
2015/04/03 12:34:56 socat[2342] N PTY is /dev/pts/6
2015/04/03 12:34:56 socat[2342] N PTY is /dev/pts/7
$ diagslave -m ascii /dev/pts/7

# Or
$ diagslave -m rtu /dev/pts/7


# RTU/ASCII Over TCP
$ socat -d -d  pty,raw,echo=0 tcp-listen:5020,reuseaddr
2018/12/25 15:57:52 socat[30337] N PTY is /dev/pts/6
2018/12/25 15:57:52 socat[30337] N listening on AF=2 0.0.0.0:5020
$ diagslave -m ascii /dev/pts/6

# Or
$ diagslave -m rtu /dev/pts/6


$ go test -v -run TCP
$ go test -v -run RTU
$ go test -v -run ASCII
$ go test -v -run RTUOverTCP
$ go test -v -run ASCIIOverTCP
```
