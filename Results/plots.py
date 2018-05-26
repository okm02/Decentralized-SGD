import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


def plot_loss(x,y,legend,title,xlabel,ylabel,picName):

	for i in range(0,len(x)):
		plt.plot(x[i],y[i])

	plt.xlabel(xlabel)
	plt.ylabel(ylabel)
	plt.title(title)
	plt.legend(legend, loc='upper left')
	plt.savefig(picName)

def readFile(fileName):
	
	with open(fileName) as f:
		content = f.readlines()

	content = [x.strip() for x in content]
	x = []
	y = []
	for i in range(0,len(content)):
		values = content[i].split(",")
		x.append(int(values[0]))
		y.append(float(values[1]))
	
	return x,y



def barPlot(title,label,data4,data5,picName):
	N = 4
	
	ind = np.arange(N)  # the x locations for the groups
	width = 0.35       # the width of the bars

	fig, ax = plt.subplots()
	rects1 = ax.bar(ind, data4, width, color='r')

	women_means = (25, 32, 34, 20, 25)
	women_std = (3, 5, 2, 3, 3)
	rects2 = ax.bar(ind + width,data5, width, color='y')

	# add some text for labels, title and axes ticks
	ax.set_ylabel(label)
	ax.set_title(title)
	ax.set_xticks(ind + width / 2)
	ax.set_xticklabels(('Ring', 'Star', 'Wheel', 'Full'))

	ax.legend((rects1[0], rects2[0]), ('4 Nodes', '5 Nodes'))
	plt.savefig(picName)

def parseFile(fileName):
	
	with open(fileName) as f:
		content = f.readlines()

	content = [x.strip() for x in content]
	data_4 = []
	data_5 = []
	for i in range(0,len(content)):
		values = content[i].split(",")
		machine = int(values[1])
		if machine == 4:
			data_4.append(float(values[2]))
		elif machine == 5:
			data_5.append(float(values[2]))
		
	
	return data_4,data_5


full4_x,full4_y = readFile('AggFull4.txt')
ring4_x,ring4_y = readFile('AggRing4.txt')
star4_x,star4_y = readFile('AggStar4.txt')
wheel4_x,wheel4_y = readFile('AggWheel4.txt')
all_x_axis = [full4_x,ring4_x,star4_x,wheel4_x]
all_y_axis = [full4_y,ring4_y,star4_y,wheel4_y]
legend = ['Fully connected','Ring','Star','Wheel']
plot_loss(all_x_axis,all_y_axis,legend,'Training Loss over iterations of dsgd','Iteration','Logistic loss','Train_Loss_4.png')


full5_x,full5_y = readFile('AggFull5.txt')
ring5_x,ring5_y = readFile('AggRing5.txt')
star5_x,star5_y = readFile('AggStar5.txt')
wheel5_x,wheel5_y = readFile('AggWheel5.txt')
all_x_axis2 = [full5_x,ring5_x,star5_x,wheel5_x]
all_y_axis2 = [full5_y,ring5_y,star5_y,wheel5_y]
legend2 = ['Fully connected','Ring','Star','Wheel']
plot_loss(all_x_axis2,all_y_axis2,legend2,'Training Loss over iterations of dsgd','Iteration','Logistic loss','Train_Loss_5.png')


exec4,exec5 = parseFile("execution.txt")
val4,val5 = parseFile("validation.txt")
barPlot('Execution time of dsgd','Execution time(sec)',exec4,exec5,"exec.png")
barPlot('Accuracy on test set of dsgd','Accuracy',val4,val5,"acc.png")







