
for i in {52..71}
do
        pm2 start populate2.sh --no-autorestart --name instance$i -- $i 
done