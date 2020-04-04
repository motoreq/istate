
for i in {26..31}
do
        pm2 start populate.sh --no-autorestart --name instance$i -- $i 
done