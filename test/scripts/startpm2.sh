
for i in {1..2}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done