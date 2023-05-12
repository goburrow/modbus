package modbus

const (
	VendorName byte = iota
	ProductCode
	MinorMajorRevision
	VendorUrl
	ProductName
	ModelNumber
	UserApplicationMame
)

type DeviceIdentification struct {
	VendorName          string
	ProductCode         string
	MinorMajorRevision  string
	VendorUrl           *string
	ProductName         *string
	ModelNumber         *string
	UserApplicationMame *string
}

func ParseDeviceIdentification(raw []byte) DeviceIdentification {
	devId := DeviceIdentification{}

	objNumber := raw[0]

	currStart := byte(1)
	for i := byte(0); i < objNumber; i++ {
		currLen := raw[currStart+byte(1)]
		currObj := raw[currStart+byte(2) : currStart+byte(2)+currLen]
		currStart += byte(2) + currLen

		strObj := string(currObj)

		switch i {
		case VendorName:
			devId.VendorName = strObj
		case ProductCode:
			devId.ProductCode = strObj
		case MinorMajorRevision:
			devId.MinorMajorRevision = strObj
		case VendorUrl:
			devId.VendorUrl = &strObj
		case ProductName:
			devId.ProductName = &strObj
		case ModelNumber:
			devId.ModelNumber = &strObj
		case UserApplicationMame:
			devId.UserApplicationMame = &strObj
		}
	}

	return devId
}
