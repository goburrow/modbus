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
$ diagslave -m ascii /dev/pts/X
```
