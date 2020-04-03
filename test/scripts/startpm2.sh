
for i in {11..15}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done