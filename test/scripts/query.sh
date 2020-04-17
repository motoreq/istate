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


        res=$(curl -sd '{"func":"QueryState","args": {"queryString":"[{\"docType\":\"eq bleh\"}]"},"user":"admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/query)
        # {"func":"QueryState","args": {"queryString":"[{\"anInt\":\"eq 1000\"}]"},"user":"admin"}

        echo $res

}

start=$(( $1 * $perWorker ))
end=$(( $start + $perWorker ))
startTime=$( date +'(%s)' )
for (( i=$start; i<=$end; i++ ))
do
	main $i
        # sleep 1s
done
endTime=$( date +'(%s)' )
echo $startTime
echo $endTime