import pandas as pd
import numpy as np
import os 
from sklearn.preprocessing import StandardScaler
from sklearn.model_selection import train_test_split

def standardize(data):

	labels = data.iloc[:,-1]
	data = data.iloc[:, :-1]
	scaler = StandardScaler()
	scaler.fit(data)
	rescaled = scaler.transform(data)
	result = pd.DataFrame(rescaled)
	result['label'] = labels
	train, test = train_test_split(result, test_size=0.2)
	train.to_csv('forest_train.csv',header=False,index=False)
	test.to_csv('forest_test.csv',header=False,index=False)
	


def checkDistribution(data):

	labels = data.iloc[:,-1]
	print(labels.value_counts())

def parseData(fname):

	with open(fname) as f:
		content = f.readlines()

	matrix = np.zeros((len(content),55))
	for i in range(0,len(content)):
		
		tokens = content[i].split()
		for j in range(1,len(tokens)):
			colVal = tokens[j].split(':')
			index = int(colVal[0])			
			val = float(colVal[1])
			matrix[i][index] = val 

		label = int(tokens[0])
		if label == 1:
			matrix[i][54] = 0
		else:
			matrix[i][54] = 1

	result = pd.DataFrame(matrix)
	return result


fileName = 'covtype.libsvm.binary'
targetFile = 'forest_train.csv'
#if not os.path.isfile(targetFile):
data = parseData(fileName)
checkDistribution(data)
#	standardize(data)


