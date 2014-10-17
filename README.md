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

References
----------
-   [Modbus Specifications and Implementation Guides](http://www.modbus.org/specs.php)
