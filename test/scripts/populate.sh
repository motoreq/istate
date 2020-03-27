#host="52.14.31.90"
host="localhost"
port="5000"
ruleID="rule_0936297c7d688faea9a6da49753f50f7905f0b81b0f09739865265a28fbcda75"
perWorker=3000
user="average"

function main() {

        ran1=$1


        res=$(curl -sd '{"func": "CreateACLRecord","args":{"resourceID":"emr_prasanths96_'$ran1'", "rulesMap":{"'$user'":{"'$ruleID'":{}}}} ,"user": "admin"}' -H "Content-Type: application/json" -X POST http://$host:$port/api/invoke)

        echo $res
}

start=$(( $1 * $perWorker ))
end=$(( $start + $perWorker ))

for (( i=$start; i<=$end; i++ ))
do
	main $i
done
