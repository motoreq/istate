
for i in {61..80}
do
        pm2 start populate.sh --name instance$i -- $i 
done