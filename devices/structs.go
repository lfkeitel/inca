package devices

type Device struct {
	Deviceid     int
	Name         string
	Hostname     string
	ConnProfile  int
	Manufacturer string
	Model        string
	Custom       bool
	Disabled     bool
	ParseConfig  bool
	Status       dStatus
}

type dStatus struct {
	Status     int
	LastPolled int
	LastError  string
}

type ConnProfile struct {
	Profileid int
	Name      string
	Protocol  string
	Username  string
	Password  string
	Enable    string
}

type DeviceStatus struct {
	Total       int
	Down        int
	Up          int
	Unknown     int
	DownDevices []Device
}

var supportedDeviceTypes = map[string][]string{
	"cisco": []string{
		"2950",
	},
	"juniper": []string{
		"2200",
	},
}

var parsableDeviceTypes = map[string][]string{
	"cisco": []string{
		"2950",
	},
	"juniper": []string{
		"2200",
	},
}

func convertIntSliceToInterface(s []int) []interface{} {
	is := make([]interface{}, len(s))
	for i, d := range s {
		is[i] = d
	}
	return is
}
