package conf

var Conf = map[string]interface{}{
	"innerAddr":  "127.0.0.1:8011",
	"outterAddr": "127.0.0.1:8088",
	"logPath":    "./log/error.log",
	"maxConn":    1000,
}
