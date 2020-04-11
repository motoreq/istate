
for i in {1..10}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done