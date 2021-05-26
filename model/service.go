package model

type Service struct {
	Name string
	Url  string
}

// GetServiceByName get service information
func GetServiceByName(serviceName string) *Service {
	rawURL := ""
	switch serviceName {
	case "test1":
		rawURL = "http://10.0.6.1:5000/"
		break
	case "test2":
		rawURL = "http://172.17.0.3:3000/"
		break
	case "debug":
		rawURL = "http://10.10.10.10/"
		break
	case "baidu":
		rawURL = "http://www.baidu.com/"
		break
	}

	return &Service{
		Name: serviceName,
		Url:  rawURL,
	}
}
