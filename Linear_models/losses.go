package Linear_models

import ("dsgd/Math"
	"math"
)

// mean squared eror loss function
func MSE(target,prediction Math.Matrix) (float64){


	batch,_ := prediction.Shape()

	diff := target.Sub(prediction)
	
	squared_diff := diff.Apply(func(a float64) float64 {return math.Pow(a,2) })

	loss := squared_diff.Sum_along_axis(0)

	return (loss.Data[0][0] / (2 * batch) )
}

// logistic regression error loss function
// We use a special form of the algorithm
// Sum xn(sigmoid(xnTw) -yn)
func Logistic(target,prediction Math.Matrix) (float64){

	logarithm := prediction.Apply(func(a float64) float64 {return math.Log(1 + math.Exp(a)) })	

	loss := logarithm.Sub(target.Dot(prediction)) 	
	
	return loss.Data[0][0]
}



// derivative of mse = 2(yn - xnTw)
func DMSE(target,prediction ,dataPoint Math.Matrix) (Math.Matrix) {
	
	batch,_ := prediction.Shape()

	gradient := (dataPoint.T()).Dot( target.Sub(prediction) )

	gradient = gradient.Apply(func(a float64) float64 {return a * (-1/batch) })
	
	return gradient
}

// derivative of logistic function
// <xn,sigmoid(xnTw) - yn)
func DLogistic(target,prediction,dataPoint Math.Matrix) (Math.Matrix){

	sigmoid := prediction.Apply(func(a float64) float64 {return Sigmoind(a) })

	gradient := (dataPoint.T()).Dot( sigmoid.Sub(target) )

	return gradient
}


func Sigmoid(prediction float64) (float64)  {

	compute := math.Exp(prediction)/(1 + math.Exp(prediction))
	return compute
}
