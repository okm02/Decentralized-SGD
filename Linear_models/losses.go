package Linear_models

import ("dsgd/Math"
	"math"
)

// mean squared eror loss function
func MSE(target,prediction,regularizer float64) (float64){

	loss:= math.Pow(target - prediction,2) + regularizer
	return loss
}

// logistic regression error loss function
// We use a special form of the algorithm
// Sum xn(sigmoid(xnTw) -yn)
func Logistic(target,prediction,regularizer float64) (float64){

	loss := math.Log(1 + math.Exp(prediction)) - (target * prediction)
	return loss + regularizer
}

// hinge loss function
func Hinge(target,prediction,regularizer float64) (float64){

	loss:= math.Max(0,1 - (target * prediction)) + regularizer
	return loss
}

// derivative of mse = 2(yn - xnTw)
func DMSE(target,prediction float64,dataPoint Math.Vector) (Math.Vector) {

	return dataPoint.Prod(prediction - target)
}

// derivative of logistic function
// <xn,sigmoid(xnTw) - yn)
func DLogistic(target,prediction float64,dataPoint Math.Vector) (Math.Vector){

	return dataPoint.Prod(Sigmoid(prediction) - target).T()
}

// derivative of hinge loss
// if hinge loss < 0 return lambda*w
// else -xnTyn + lambda //w//^2
func DHinge(target,hinge_loss float64,dataPoint,regularizer Math.Vector) (Math.Vector) {

	if hinge_loss <= 0.0 {
		return regularizer
	}else{
		sum,_ := dataPoint.T().Prod(-1 * target).Sum(regularizer)
		return sum
	}

}

func Sigmoid(prediction float64) (float64)  {

	compute := math.Exp(prediction)/(1 + math.Exp(prediction))
	return compute
}