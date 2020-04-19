# Copyright 2020 Prasanth Sundaravelu

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
perWorker=100

function main() {

        ran1=$1


        #res=$(curl -sd '{"func":"CreateState","args": {"docType":"bleh", "id":"bleh'$ran1'", "anInt": 100000, "aMultiStruct":{"multiVal":{"val":"multivalstring"}}, "anArray":[-1,-1,-1],"a3DArray":[[[1],[2]],[[3]]], "a2DArray":[[1,2],[3,4]] ,"aMap":{"1":-5,"2":-5}, "aStruct":{"val":"astructvalue"}, "aComplexMapSlice":[{"1stindex":[{"1": {}}]}], "aMapStruct":[{"mapfield1":{"val":"asv"}}]}, "user":"admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/invoke)
        res=$(curl -sd '{"func":"CreateState","args": {"docType":"bleh", "id":"bleh'$ran1'", "anInt": 1000, "aMultiStruct":{"multiVal":{"val":"multivalstring"}}, "aStruct":{"val":"astructvalue"}}, "user":"admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/invoke)
        echo $res
}

start=$(( $1 * $perWorker ))
end=$(( $start + $perWorker ))

for (( i=$start; i<=$end; i++ ))
do
	main $i
done
