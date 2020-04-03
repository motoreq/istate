
for i in {22..26}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done