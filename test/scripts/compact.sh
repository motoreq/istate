# Copyright 2020 Motoreq Infotech Pvt Ltd

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#host="52.14.31.90"
host="localhost"
port="5000"
num=$1

function main() {
        res=$(curl -sd '{"func":"CompactIndex","args": {}, "user":"admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/invoke)
        echo $1
        echo $res
}

start=1
for (( i=$start; i<=$num; i++ ))
do
	main $i
done