package main

import (
	"os"
	"fmt"
	"net"
	"protobuf"
	"math/rand"
	"time"
	"strconv"
	"dsgd/util"
	"dsgd/Math"
	"dsgd/Messages"
	"dsgd/Linear_models"
	"sync"
)

// a struct used to move data between different threads
type Update struct {
	weight Math.Vector
	end bool
}


func main()  {

	// parse standard input
	params := os.Args[1:]
	port := util.ReadPort(params[0])
	startPeer(port)
}

func startPeer(port  net.UDPAddr){


	myAddr,addr_err := net.ResolveUDPAddr("udp4",port.String())
	util.CheckError(addr_err)
	conn,conn_err := net.ListenUDP("udp4",myAddr)
	util.CheckError(conn_err)
	defer conn.Close()

	fmt.Println("Started send me a message")
	decodeMsgTemp := &Messages.Init_Message{}
	data_buffer := make([]byte, 1024)
	_, _, err := conn.ReadFromUDP(data_buffer)
	util.CheckError(err)
	err = protobuf.Decode(data_buffer, decodeMsgTemp) // receive initial parameter message from master

	fmt.Println("Received init message successfully")


	path:= "./Data/process/train0" + strconv.Itoa(decodeMsgTemp.PID)
	data := util.ReadData(path) // read my part of the data
	final_weights := stochastic_gradient_descent(*conn,data,*decodeMsgTemp) // run sgd
	time.Sleep(2 * time.Second) // give the master time to collect all partial loss functions
	fmt.Println("Woke up was a nice sleep")
	send_FinalWeights(final_weights,decodeMsgTemp.PID,decodeMsgTemp.MasterPort) // send the final weight
	conn.Close()

}

func stochastic_gradient_descent(conn net.UDPConn,data []util.Data_point,parameters Messages.Init_Message) (Math.Vector) {

	weights := Math.Zeros(data[0].Data.Col_dimension,1) // initial weights
	var losses []util.LossT // array of losses
	channel := make(chan Update)

	for i:=0;i<parameters.Max_iterations ;i++  {

		index := random_index(len(data)) // select a random index
		current_datapoint := data[index] // get the current datapoint
		prediction,_ := current_datapoint.Data.Dot(weights) // multiply data point with weight vector

		regularizer := weights.L2Norm() * 0.0001 // compute regularization
		loss := Linear_models.Logistic(current_datapoint.Label,prediction,regularizer) // compute logistic loss
		fmt.Printf("Loss = %.6f\n", loss)

		losses = append(losses, util.LossT{i,loss})

		// main thread computes gradient,meanwhile a thread is launched to fetch weights from neighbouring
		// workers
		gradient := Linear_models.DLogistic(current_datapoint.Label,prediction,current_datapoint.Data)
		go handle_incoming_msgs(conn,channel,parameters.Stoch_row,weights.Row_dimension,weights.Col_dimension)
		request_weights(parameters.Neighbours,parameters.Stoch_row,parameters.PID,weights)
		weighted_sum := <- channel

		if weighted_sum.end { // a checkpoint to break if loss on the total dataset is achieved
			break
		}else {
			// compute the local weight contribution by multipliying my weight with the stochasic weight
			// then perform the stochastic gradient update
			local_contribution := weights.Prod(parameters.Stoch_row[parameters.PID])
			total_weighted_sum,_ := weighted_sum.weight.Sum(local_contribution)
			temp,err := total_weighted_sum.Subtract(gradient.Prod(parameters.Step_size))
			util.CheckError(err)
			weights = temp
			// send loss on my dataset to master
			send_PartialLoss(i,loss,parameters.MasterPort,false)
		}

	}
	close(channel)
	// message to denote end of computation
	send_PartialLoss(-1,-1,parameters.MasterPort,true)
	return weights

}

// send my weight vector to workers ,whom have a value on my stochastic row > 0
func request_weights(neighbours []net.UDPAddr,stochastic_row []float64,myId int,weight Math.Vector){

	for i:=0;i<len(stochastic_row);i++{

		if stochastic_row[i]>0 && i!=myId {

			peer := neighbours[i]
			Gradient_msg := &Messages.Request_Weights{weight,myId,false}
			conn,conn_err := net.Dial("udp",peer.String())
			util.CheckError(conn_err)
			packetMessage,msgErr := protobuf.Encode(Gradient_msg)
			util.CheckError(msgErr)
			_,sendErr := conn.Write(packetMessage)
			util.CheckError(sendErr)
			conn.Close()
		}

	}
}


func handle_incoming_msgs(conn net.UDPConn,sum chan <- Update,stochastic_row []float64,rowDim int,colDim int) {

	channel := make(chan Update)
	weighted_sum := Math.Zeros(rowDim,colDim)
	non_zeros := util.Non_zero(stochastic_row) - 1
	counter :=0
	computation_done := false
	var mutex = &sync.Mutex{}
	for {

		decodeMsgTemp := &Messages.Request_Weights{}
		data_buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(data_buffer)
		util.CheckError(err)
		err = protobuf.Decode(data_buffer, decodeMsgTemp)
		// read weight sent by connected node

		// launch a thread to handle the incoming message
		go handle_peers(channel,*decodeMsgTemp,stochastic_row)
		weighted_vector := <- channel
		mutex.Lock()
		if !weighted_vector.end { // if not end of computation
			// add the weighted optimization variable to the aggregated weight sums
			weighted_sum,_ =weighted_sum.Sum(weighted_vector.weight)
			counter = counter + 1
		}else {
			computation_done = true
		}
		mutex.Unlock()
		if counter >= non_zeros || computation_done{ // I received weights from all
			break
		}
	}
	close(channel)
	sum <- Update{weighted_sum,computation_done}
}


func handle_peers(channel chan <- Update ,msg Messages.Request_Weights,stochastic_row []float64){

	if msg.EndComputation {
		channel <- Update{Math.Vector{},true}
	}else {
		stochastic_value := stochastic_row[msg.PID]
		weighted_sum := msg.Weight.Prod(stochastic_value) // multiply the received vector with its stochastic weight
		channel <- Update{weighted_sum,false}
	}

}

func send_PartialLoss(iter int,loss float64,masterPort string,end bool)  {

	msg := &Messages.Request_Loss{iter,loss,end}
	conn,conn_err := net.Dial("udp",masterPort)
	util.CheckError(conn_err)
	packetMessage,msgErr := protobuf.Encode(msg)
	util.CheckError(msgErr,)
	_,sendErr := conn.Write(packetMessage)
	util.CheckError(sendErr)
	conn.Close()
}

func send_FinalWeights(Result Math.Vector,Pid  int,masterPort string){

	msg := &Messages.Request_Weights{Weight:Result,PID:Pid,EndComputation:true}
	conn,conn_err := net.Dial("udp",masterPort)
	util.CheckError(conn_err)
	packetMessage,msgErr := protobuf.Encode(msg)
	util.CheckError(msgErr)
	_,sendErr := conn.Write(packetMessage)
	util.CheckError(sendErr)
}



func random_index(points int) (int) {

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return r.Intn(points)

}

