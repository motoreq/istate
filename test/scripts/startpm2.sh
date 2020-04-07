
for i in {103..44}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done