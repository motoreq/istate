awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' ./orderer.pem > orderer.txt
awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' ./peer.pem > peer.txt