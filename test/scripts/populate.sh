#host="52.14.31.90"
host="localhost"
port="5000"
perWorker=100

function main() {

        ran1=$1


        res=$(curl -sd '{"func":"CreateState","args": {"docType":"bleh", "id":"bleh'$ran1'", "anInt": 600, "aMultiStruct":{"multiVal":{"val":"multivalstring"}}, "anArray":[-1,-1,-1],"a3DArray":[[[1],[2]],[[3]]], "a2DArray":[[1,2],[3,4]] ,"aMap":{"1":-5,"2":-5}, "aStruct":{"val":"astructvalue"}, "aComplexMapSlice":[{"1stindex":[{"1": {}}]}], "aMapStruct":[{"mapfield1":{"val":"asv"}}]}, "user":"admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/invoke)
        echo $res
}

start=$(( $1 * $perWorker ))
end=$(( $start + $perWorker ))

for (( i=$start; i<=$end; i++ ))
do
	main $i
done
