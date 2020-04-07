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