package util

import (
	"strconv"
	"log"
	"net"
	"strings"
	"os"
	"bufio"
	"dsgd/Math"
	"fmt"
)

// data point encapsulates the data features stored in the data vector
// label correpsonds to the class of the datapoint
type Data_point struct {
	Data  Math.Vector
	Label float64
}

// a tuple used to write output of each node to the file
type LossT struct {
	Iteration int
	Loss float64
}

// parse a port string into network address
func ReadPort (input string) (net.UDPAddr) {

	getPorts := strings.Split(strings.Split(input,":")[1],",")

	localHost := "127.0.0.1:" + getPorts[0]
	addr,_ := net.ResolveUDPAddr("udp",localHost)

	return *addr
}

// create a zero matrix used for creation of link matrix
func zero_matrix(row_dim,col_dim int)  ([][]float64) {

	mat := make([][]float64,row_dim)
	for i:=0;i<row_dim ;i++  {
		mat[i] = make([]float64,col_dim)
	}
	return mat
}

// sum over a row of a matrix
func axis_sum(matrix []float64)  (float64){

	sum :=0.0
	for i:=0;i<len(matrix);i++  {
		sum+= matrix[i]
	}
	return sum
}


// get a link matrix representing the link matrix
// normalize each row
func Stochastic_matrix(fileName string,numPIS int) ([][] float64) {

	linkMatrix := readTopology(fileName,numPIS)
	for i:=0;i< numPIS;i++ {
		rowSum := axis_sum(linkMatrix[i])
		normalized := make([]float64,len(linkMatrix[i]))
		for j:=0;j<len(linkMatrix[i]);j++  {
			normalized[j] = linkMatrix[i][j]/rowSum
		}
		linkMatrix[i] = normalized
	}
	return linkMatrix
}

// read the topology from a text file
// input is of the form nodeID1 nodeID2
// coresponding to the existence of a link between those two nodes
func readTopology(fileName string,numPIS int) ([][] float64) {

	linkMatrix  := zero_matrix(numPIS,numPIS)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		nodes := strings.Split(scanner.Text()," ")

		outgoing,_ := strconv.Atoi(nodes[0])
		incoming,_ := strconv.Atoi(nodes[1])
		linkMatrix[outgoing][incoming] = 1.0
		linkMatrix[incoming][outgoing] = 1.0
	}
	file.Close()

	for i:=0;i<len(linkMatrix);i++{
		linkMatrix[i][i] = 1.0
	}

	return linkMatrix
}


// parse a list of ip ports
func Parse_Neighbours(input string) (NeighbourAddress []net.UDPAddr) {

	getPorts := strings.Split(strings.Split(input,":")[1],",")

	var peers []net.UDPAddr

	for i:=0;i<len(getPorts);i++ {
		localHost := "127.0.0.1:" + getPorts[i]
		addr,_ := net.ResolveUDPAddr("udp",localHost)
		peers = append(peers,*addr)
	}

	return peers
}

// function to read the data files
// @output : datapoint array storing features and labels of each row in the data matrix
func ReadData(dataPath string) ([]Data_point) {

	var representation []Data_point

	inputD,erra := os.Open(dataPath)
	if erra!=nil  {
		log.Fatal("Error loading files")
	}
	defer inputD.Close()
	scanner := bufio.NewScanner(inputD)
	rowId := 0
	for scanner.Scan(){
		label,data := parseRow(scanner.Text())
		temp := Data_point{Math.Vector{1,
			len(data),
			data},label}
		representation = append(representation,temp)
		rowId++
	}
	inputD.Close()

	return representation
}

// function used in read data function to parse the text files
func parseRow(input string)  (float64,[]float64) {

	parsing := strings.Split(input,",")
	dimension := len(parsing)
	label,_ := strconv.ParseFloat(parsing[dimension - 1],64)
	data := make([]float64,dimension-1)
	for i:=1;i<dimension-1;i++{
		parsed,_ := strconv.ParseFloat(parsing[i],64)
		data[i] = parsed
	}
	return label,data
}

// count number of nonzero inputs
// used to check how many messages a process is expecting from his peers
func Non_zero(row []float64) int {
	count := 0
	for i:=0;i<len(row) ;i++  {
		if row[i]>0.0 {
			count+=1
		}
	}
	return count
}



func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// write loss to file
func OutputLoss(fileName string,loss []LossT)  {

	file , err := os.Create("./Results/" + fileName + ".txt")
	CheckError(err)
	defer file.Close()
	for i:=0;i<len(loss) ;i++ {
		file.WriteString(fmt.Sprintf("%d,%f\n",loss[i].Iteration,loss[i].Loss))
	}
	file.Close()
}