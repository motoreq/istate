
#!/bin/bash
docker rmi -f `docker images dev* -aq`
