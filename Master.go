package main

import (
	"net"
	"dsgd/util"
	"dsgd/Messages"
	"protobuf"
	"sync"
	"os"
	"strconv"
	"strings"
	"dsgd/Math"
	//"math"
	"time"
	"dsgd/Linear_models"
)

func main()  {
	master_computation()
}

func master_computation() {

	// parse standard input
	peerIds := os.Args[1]
	masterPort := util.ReadPort(os.Args[2])
	numPIS, _ := strconv.Atoi((strings.Split(os.Args[3], ":"))[1])
	topo := strings.Split(os.Args[4], ":")[1]
	peers := util.Parse_Neighbours(peerIds)
	step_size := 0.01
	max_iters := 10000

	topology := "./Topology/" + topo + strconv.Itoa(numPIS) + ".txt"
	link := util.Stochastic_matrix(topology, numPIS)

	// start computation
	startTime := time.Now()
	// send messages to workers to start the computation
	start_computation(link, peers, step_size, max_iters, numPIS, masterPort.String(), topo)
	// bind to network port to aggregate loss functions
	myAddr, addr_err := net.ResolveUDPAddr("udp4", masterPort.String())
	util.CheckError(addr_err)
	conn, conn_err := net.ListenUDP("udp4", myAddr)
	util.CheckError(conn_err)
	defer conn.Close()
	// aggregate losses at every iteration
	agg_loss(*conn, peers, 0.001, topo)

	channel := make(chan Math.Vector)
	// compute the sum of all weight vectors at the last iteration of the computation
	go acummilate_weights(*conn, numPIS, channel)
	weight := <-channel
	close(channel)
	conn.Close()
	done := time.Since(startTime)
	// compute loss on test data
	validation_set := "./Data/process/forest_test.csv"
	normalized_weight := weight.Div(float64(numPIS)) // total weight is the averaged sum of weights
	accuracy := strconv.FormatFloat(compute_validation_loss(normalized_weight, validation_set), 'E', -3, 64)
	outputFile("./Results/execution.txt", done.String(), topo, numPIS)
	outputFile("./Results/validation.txt",accuracy,topo,numPIS)
}

func start_computation(link [][]float64,peers []net.UDPAddr,step_size float64,max_iter int,pis int,masterPort string,topo string){

	for i:=0;i<pis;i++{
		// send initialization messages to each worker node
		Start_msg := Messages.Init_Message{i,masterPort,link[i],peers,step_size,max_iter,topo}
		conn,conn_err := net.Dial("udp",peers[i].String())
		util.CheckError(conn_err)
		packetMessage,msgErr := protobuf.Encode(&Start_msg)
		util.CheckError(msgErr)
		_,sendErr := conn.Write(packetMessage)
		util.CheckError(sendErr)
		conn.Close()

	}

}

func agg_loss(conn net.UDPConn,peers[]net.UDPAddr,epsilon float64,topology string)  {

	mutex:= &sync.Mutex{}
	//prev_loss := 0.0
	curr_loss := 0.0
	iteration := 0
	count := 0
	//delta := epsilon
	channel := make(chan float64)
	condition := false
	var losses []util.LossT
	for {
		// handle incoming messages having loss of each worker at each iteration
		decodeMsgTemp := &Messages.Request_Loss{}
		data_buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(data_buffer)
		util.CheckError(err)
		err = protobuf.Decode(data_buffer, decodeMsgTemp)

		// after decoding the message launch a thread to handle it
		go handle_incoming_loss(*decodeMsgTemp,channel)
		partial_loss := <- channel

		mutex.Lock()

		if partial_loss>0 {
			curr_loss = curr_loss + partial_loss // aggregate loss
			count = count + 1
		} else {
			condition = true // denotes end of computation
		}


		if count == len(peers) {
			 count = 0
			 iteration = iteration + 1
			 // store aggregated losses to save to output
			 losses = append(losses,util.LossT{iteration,curr_loss})
			 //delta = math.Abs(curr_loss - prev_loss)
			// prev_loss = curr_loss
			 curr_loss = 0.0 // reset loss for next computation
		}
		mutex.Unlock()
		if condition{
			break
		}
		/*if delta < epsilon {

			for i:=0;i<len(peers) ; i++  {
				End_msg := &Messages.Request_Weights{Math.Vector{},-1,true}
				conn,conn_err := net.Dial("udp",peers[i].String())
				util.CheckError(conn_err,"Master connection error")
				packetMessage,_ := protobuf.Encode(End_msg)
				_,sendErr := conn.Write(packetMessage)
				util.CheckError(sendErr,"Master message sending error")
				conn.Close()
			}
			fmt.Printf("I converged at iteration =  %d \n",iteration)
			break

		}*/
	}
	close(channel)
	util.OutputLoss("Agg" + topology + strconv.Itoa(len(peers)),losses) // write total loss
}


func acummilate_weights(conn net.UDPConn,pis int,channel chan <-Math.Vector) {


	tempChann := make(chan Math.Vector)
	var weights []Math.Vector
	counter := 0
	mutex:= &sync.Mutex{}
	for {
		// read final weight messages
		decodeMsgTemp := &Messages.Request_Weights{}
		data_buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(data_buffer)
		util.CheckError(err)
		err = protobuf.Decode(data_buffer, decodeMsgTemp)
		if err == nil { // some messages loss messages may come late hence we ignore them
			go handle_incoming_weight(*decodeMsgTemp,tempChann)
			weightPI := <- tempChann
			mutex.Lock()
			weights = append(weights,weightPI) // append all losses to a list of vectors
			counter = counter + 1
			mutex.Unlock()
			if counter>= pis {
				break
			}
		}

	}
	close(tempChann)
	total := Math.Zeros(weights[0].Row_dimension,weights[0].Col_dimension)
	// compute sum of final weight vectors
	for i:=0;i<len(weights) ;i++  {
		total,_ = total.Sum(weights[i])
	}
	channel <- total
}

func handle_incoming_loss(msg Messages.Request_Loss,channel chan <- float64){
	channel <- msg.Loss
}

func handle_incoming_weight(msg Messages.Request_Weights,channel chan <- Math.Vector){
	channel <- msg.Weight
}

func compute_validation_loss(weight Math.Vector,fileName string) float64 {
	validation_set := util.ReadData(fileName) // read the test set
	validation_size := float64(len(validation_set)) // compute its size
	correct := 0.0
	mistake := 0.0
	for i:=0;i<len(validation_set) ;i++  {
		prediction,_ := validation_set[i].Data.Dot(weight) // prediction is xnTw
		prediction_class := quantize_probabilities(Linear_models.Sigmoid(prediction)) // compute label prediction
		if prediction_class == validation_set[i].Label { // if same class increase correct predictions
			correct = correct + 1.0
		}else {
			mistake = mistake + 1
		}
	}
	accuracy := correct/validation_size // compute accuracy
	return accuracy
}

// function to quantize probabilites generated by sigmoid function
func quantize_probabilities(proba float64) float64  {

	if proba >=0.5 {
		return 1.0
	}else {
		return 0.0
	}

}

// outputs data to append only files
func outputFile(filename string,execTime string,topology string,pis int){
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	text := topology + "," + strconv.Itoa(pis) + "," + execTime + "\n"
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
	f.Close()
}