go modbus
=========
Fault-tolerant, fail-fast implementation of Modbus protocol in Go.

Incubating - do not use in production yet (more testing needed).

Supported functions
-------------------
Bit access:
*   Read Discrete Inputs
*   Read Coils
*   Write Single Coil
*   Write Multiple Coils

16-bit access:
*   Read Input Registers
*   Read Holding Registers
*   Write Single Register
*   Write Multiple Registers
*   Read/Write Multiple Registers
*   Mask Write Register
*   Read FIFO Queue

Supported formats
-----------------
*   TCP
*   Serial (RTU, ASCII) - in progress

Usage
-----
Basic usage:
```go
client := modbus.TcpClient("localhost:502")
// Read input register 9
results, err := client.ReadInputRegisters(8, 1)
```

Advanced usage:
```go
var handler modbus.TcpClientHandler
handler.ConnectString = "localhost:502"
handler.Timeout = 5 * time.Second
handler.UnitId = 0x01
handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
// Connect manually
handler.Connect()
defer handler.Close()

client := modbus.TcpClientWithHandler(&handler)
results, err := client.ReadDiscreteInputs(15, 2)
results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
```

References
----------
-   [Modbus Specifications and Implementation Guides](http://www.modbus.org/specs.php)
