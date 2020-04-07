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