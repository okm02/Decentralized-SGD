package Messages

import (
	"net"
	"dsgd/Math"
)

// initialization method
// used to communicate parameters of the model to start computation
type Init_Message struct {

	PID int
	MasterPort string
	Stoch_row []float64
	Neighbours []net.UDPAddr
	Step_size float64
	Max_iterations int
	Topology string
}

// message used to exchange weights among peers
type Request_Weights struct {
	Weight Math.Vector `json:"weight"`
	PID int				`json:"pid"`
	EndComputation bool	`json:"end_computation"`
}

// a message used by master to aggregate local loss at every node
type Request_Loss struct {
	Iteration int
	Loss float64
	EndComputation bool
}



