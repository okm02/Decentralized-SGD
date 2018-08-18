package Math

import (
	"log"
)

/*
 * change matrix representation such that array level1 corresponds to columns and level 2 to rows
 */

type Matrix struct {
	row_dimension int
	col_dimension int
	Data [][] float64
}

// constructor
func New(row_dim,col_dim int) (Matrix){

	return Matrix{
		row_dimension:row_dim,
		col_dimension:col_dim,
		Data:zeros(row_dim,col_dim),
	}

}

func (this Matrix) Shape() (int,int){
	return this.row_dimension,this.col_dimension
}

func zeros(row_dim,col_dim int) ([][]float64){

	zero_matrix := make([][]float64,row_dim)
	for i:=0;i<len(zero_matrix);i++{

		zero_matrix[i] = make([]float64,col_dim)
		for j:=0;j<len(zero_matrix[i]) ;j++  {
			zero_matrix[i][j] = 0.0
		}
	}
	return zero_matrix
}


func (this Matrix) T() (Matrix){

	transposed_matrix := make([][] float64,this.col_dimension)
	for i:=0;i<this.col_dimension;i++ {

		transposed_matrix[i] = make([]float64,this.row_dimension)
		for j:=0;j<this.row_dimension;j++{
			transposed_matrix[i][j] = this.Data[j][i]
		}
	}

	return Matrix{this.col_dimension,this.row_dimension,transposed_matrix}

}

func (this Matrix) Dot (other Matrix) (Matrix){

	vector_assertion := (other.col_dimension == 1)
	dimension_assertion := (this.col_dimension == other.row_dimension)
	if vector_assertion && dimension_assertion {

		result := New(this.row_dimension,1)

		for i:=0;i<this.row_dimension;i++{
			
			sum := 0.0
			for j:=0;j<this.col_dimension;j++ {
				sum+= this.Data[i][j] * other.Data[j][0]
			}
			result.Data[i][0] = sum
		}

		return result
	}else{
		log.Fatal("Operation non supported")
	}

	return Matrix{}

}


func (this Matrix) Apply ( op func(a float64) (float64) ) (Matrix) {

	for i:=0;i<this.row_dimension;i++{
		for j:=0;j<this.col_dimension;j++ {
			this.Data[i][j] = op(this.Data[i][j])
		}
	}
	return this
}

func element_wise_ops(matrixA [][]float64,matrixB [][]float64,op func(a,b float64) (float64))  ([][]float64) {
	
	// performs elementwise operations between two matrices

	row_dimension := len(matrixA)
	col_dimension := len(matrixA[0])
	temp_data := make([][]float64,row_dimension)

	for i:=0;i<row_dimension;i++{

		temp_data[i] = make([]float64,col_dimension)

		for j:=0;j<col_dimension;j++{

			temp_data[i][j] = op(matrixB[i][j], matrixA[i][j])
		}
	}

	//fmt.Println(temp_data)

	return temp_data
}



func (this Matrix) Add(other Matrix) (Matrix)  {

	dim_assertion := (this.row_dimension == other.row_dimension) && (this.col_dimension == other.col_dimension)
	if dim_assertion{

		summed_matrix := element_wise_ops(this.Data,other.Data, func(a, b float64) float64 {return a+b} )		

		return Matrix{this.row_dimension,this.col_dimension,summed_matrix}
	}else{
		log.Fatal("Non equivalent dimensions")
	}

	return Matrix{}
}


func (this Matrix) Sub(other Matrix) (Matrix)  {

	dim_assertion := (this.row_dimension == other.row_dimension) && (this.col_dimension == other.col_dimension)
	if dim_assertion{

		summed_matrix := element_wise_ops(this.Data,other.Data, func(a, b float64) float64 {return a-b} )		

		return Matrix{this.row_dimension,this.col_dimension,summed_matrix}
	}else{
		log.Fatal("Non equivalent dimensions")
	}

	return Matrix{}
}

func (this Matrix) Mul(other Matrix) (Matrix)  {

	dim_assertion := (this.row_dimension == other.row_dimension) && (this.col_dimension == other.col_dimension)
	if dim_assertion{

		summed_matrix := element_wise_ops(this.Data,other.Data, func(a, b float64) float64 {return a*b} )		

		return Matrix{this.row_dimension,this.col_dimension,summed_matrix}
	}else{
		log.Fatal("Non equivalent dimensions")
	}

	return Matrix{}
}

func (this Matrix) Div(other Matrix) (Matrix)  {

	dim_assertion := (this.row_dimension == other.row_dimension) && (this.col_dimension == other.col_dimension)
	if dim_assertion{

		summed_matrix := element_wise_ops(this.Data,other.Data, func(a, b float64) float64 {return a/b} )		

		return Matrix{this.row_dimension,this.col_dimension,summed_matrix}
	}else{
		log.Fatal("Non equivalent dimensions")
	}

	return Matrix{}
}


func (this Matrix) Sum_along_axis(axis int) (Matrix){


	axis_assertion := (axis == 0) || (axis ==1)
	
	if axis_assertion{

		if axis == 0{

			row_vector := New(this.row_dimension,1)
			for i:=0;i<this.row_dimension;i++{
				sum := 0.0
				for j:=0;j<this.col_dimension;j++{
					sum+= this.Data[i][j]
				}
				row_vector.Data[i][0] = sum
			}
			
			return row_vector

		}else{

			col_vector := New(1,this.col_dimension)
			for i:=0;i<this.col_dimension;i++{
				sum := 0.0
				for j:=0;j<this.row_dimension;j++{
					sum+= this.Data[j][i]
				}
				col_vector.Data[i][0] = sum
			}
			
			return col_vector

		}

	}else{
		log.Fatal("Axis not supported")
	}

	return Matrix{}

}








