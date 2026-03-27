run:
	DEBUG_ENABLED=x INTERVAL=30 DEVICES=energysocket1:192.168.0.10,energysocket2:192.168.0.11,p1-meter:192.168.0.12,watermeter:192.168.0.13 go run cmd/main.go
