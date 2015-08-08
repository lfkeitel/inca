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
	Status       dStatus
}

type dStatus struct {
	Status     int
	LastPolled int
}

type ConnProfile struct {
	Profileid int
	Name      string
	Protocol  string
	Username  string
	Password  string
	Enable    string
}
