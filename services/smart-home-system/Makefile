SENSOR=sensor1
current:
	@ curl 127.0.0.1:7070/sensor/$(SENSOR)/current-temperature

VALUE=123
set:
	@ curl --request POST 127.0.0.1:7070/sensor/$(SENSOR)/target-temperature/$(VALUE)