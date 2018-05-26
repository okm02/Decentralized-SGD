package Math

import (
	"errors"
	"math"
)

// a vector is constructed having row and columns of the input
// data contains our points
type Vector struct {

	Row_dimension int `json:"row_dimension"`
	Col_dimension int `json:"col_dimension"`
	Data []float64	   `json:"data"`
}

// check if the row and column dimensions of the two vectors are equivalent
func (row_vec Vector) Dot(col_vec Vector) (sum float64,err error){

	if row_vec.Col_dimension != col_vec.Row_dimension {
		return 0,errors.New("Incompatible dimensions")
	}

	for i:=0;i<len(row_vec.Data);i++{
		sum+= (row_vec.Data[i] * col_vec.Data[i])
	}
	return
}

// check for equivalent dimensions
// then sum the data arrays
func (row_vecA Vector) Sum (row_vecB Vector) (sum Vector,err error){

	if row_vecA.Row_dimension!= row_vecB.Row_dimension || row_vecA.Col_dimension!=row_vecB.Col_dimension {
		return Vector{},errors.New("Incompatible dimensions")
	}

	result := make([]float64,len(row_vecA.Data))
	for i:=0;i<len(row_vecA.Data);i++{
		result[i] = (row_vecA.Data[i] + row_vecB.Data[i])
	}
	sum.Row_dimension = row_vecA.Row_dimension
	sum.Col_dimension = row_vecB.Col_dimension
	sum.Data = result
	return
}

// same as addition function
func (row_vecA Vector) Subtract (row_vecB Vector) (sum Vector,err error){

	if row_vecA.Row_dimension!= row_vecB.Row_dimension || row_vecA.Col_dimension!=row_vecB.Col_dimension {
		return Vector{},errors.New("Incompatible dimensions")
	}

	result := make([]float64,len(row_vecA.Data))
	for i:=0;i<len(row_vecA.Data);i++{
		result[i] = (row_vecA.Data[i] - row_vecB.Data[i])
	}
	sum.Row_dimension = row_vecA.Row_dimension
	sum.Col_dimension = row_vecB.Col_dimension
	sum.Data = result
	return
}

// transpose in this simplistic schema is to switch the dimensions
func (vec Vector) T () Vector{

	temp_col := vec.Col_dimension
	temp_row := vec.Row_dimension
	vec.Col_dimension = temp_row
	vec.Row_dimension = temp_col
	return vec

}

// l2 norm of a vector is sum of each input squared
// then take square root
func (vec Vector) L2Norm () float64  {

	squared_sum := 0.0
	for i:=0;i < len(vec.Data) ; i++  {
		squared_sum += (math.Pow(vec.Data[i],2))
	}
	return math.Sqrt(squared_sum)
}

// create a zero row vector
func  Zeros (row_dim,col_dim int) Vector  {

	res := make([]float64,row_dim)
	for i:=0;i<len(res);i++   {
		res[i] = 0.0
	}

	return Vector{
		Row_dimension:row_dim,
		Col_dimension:col_dim,
		Data:res,
	}
}

/*
* scalar-vector ops
 */

func (vec Vector) Add (scalar float64) Vector  {

	temp := make([]float64,len(vec.Data))
	for i:=0;i<len(temp) ; i++  {
		temp[i] = vec.Data[i] + scalar
	}
	vec.Data = temp
	return vec
}

func (vec Vector) Diff (scalar float64) Vector  {

	temp := make([]float64,len(vec.Data))
	for i:=0;i<len(temp) ; i++  {
		temp[i] = vec.Data[i] - scalar
	}
	vec.Data = temp
	return vec
}

func (vec Vector) Prod (scalar float64) Vector  {

	temp := make([]float64,len(vec.Data))
	for i:=0;i<len(temp) ; i++  {
		temp[i] = vec.Data[i] * scalar
	}
	vec.Data = temp
	return vec
}

func (vec Vector) Div (scalar float64) Vector  {

	temp := make([]float64,len(vec.Data))
	for i:=0;i<len(temp) ; i++  {
		temp[i] = vec.Data[i] / scalar
	}
	vec.Data = temp
	return vec
}