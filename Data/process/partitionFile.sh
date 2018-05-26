#!/bin/bash
# Varaibles containing the file and number of processes

fspec=forest_train.csv
num_files=5

# number of lines per file

total_lines=$(wc -l <${fspec})
((lines_per_file = (total_lines + num_files - 1) / num_files))

# Split the actual file, maintaining lines.

split -d --lines=${lines_per_file} ${fspec} train

# Debug information

echo "Total lines     = ${total_lines}"
echo "Lines  per file = ${lines_per_file}"    
#wc -l train.*
