# Decentralized Stochastic Gradient Descent



## Prerequistes

To run this project please install the go programming language on your machine and set the gopath variable.

## DATA
Dataset is under the following link : [data](https://www.csie.ntu.edu.tw/~cjlin/libsvmtools/datasets/binary.html) and it's called : covtype.binary

## Running the tests

To run the code you need to execute the lines below.

```
cd ./Topology
python topologies.py 5
cd ../Data/process
python dataProcessing.py
chmod +x partitionFile.sh
./partitionFile.sh
cd ..
make target
```

Just some warning if you need to add or decrese the number of workers. You need to execute the topologies script with the number of processes you need as standard input. Also you need to change the num_files in the partitionFile script. Finally in the makefile you should add do the folllowing :
```
./Peer -port:newPort &
sleep1
```
For the master :
```
./Master -peers:8000,8001,8002,8003,8004,newPort -master:6000 -pis:newNumberOfPIS -topo:Full
``` 
To set the topology you want to run, please set the topo field in master to the one you want.


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## References
- Lian, Xiangru; Zhang, Ce; Zhang et al Can Decentralized Algorithms Outperform Centralized Algorithms? A Case Study for Decentralized Parallel Stochastic Gradient Descent. In ArXiv e-prints.
