import networkx as nx
import sys

def fully_connected(nodes,path):

	graph = nx.complete_graph(nodes)
	edges = [e for e in graph.edges()]
	write_topology(nodes,edges,'Full',path)
	

def star(nodes,path):
	graph = nx.star_graph(nodes)
	edges = [e for e in graph.edges()]
	write_topology(nodes+1,edges,'Star',path)

def ring(nodes,path):
	graph = nx.cycle_graph(nodes)
	edges = [e for e in graph.edges()]
	write_topology(nodes,edges,'Ring',path)

def wheel(nodes,path):
	graph = nx.wheel_graph(nodes)
	edges = [e for e in graph.edges()]
	write_topology(nodes,edges,'Wheel',path)

def random(nodes,path):
	graph = nx.gnm_random_graph(nodes,nodes)
	edges = [e for e in graph.edges()]
	write_topology(nodes,edges,'Random',path)

def write_topology(nodes,edges,topoName,path):

	fileName = path + topoName + str(nodes) + ".txt"
	f = open(fileName,'w') 
	for i in edges:
		f.write(str(i[0]) + " " + str(i[1]) + "\n")
	f.close()

def generate_and_write(nodes,path):

	fully_connected(nodes,path)
	star(nodes-1,path)
	ring(nodes,path)
	wheel(nodes,path)
	random(nodes,path)


nodes = sys.argv[1]
generate_and_write(int(nodes),'./')
