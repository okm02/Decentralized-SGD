Topo =	Ring\
	Star \
	Wheel \
	Full 


target: Master Peer									
	./Peer -port:8000 &
	sleep 1
	./Peer -port:8001 &
	sleep 1
	./Peer -port:8002 &
	sleep 1
	./Peer -port:8003 &
	sleep 1		  &
	./Peer -port:8004 &
	sleep 1		  &	
	./Master -peers:8000,8001,8002,8003,8004 -master:6000 -pis:5 -topo:Full

Master:Master.go
	go build Master.go

Peer:Peer.go
	go build Peer.go
