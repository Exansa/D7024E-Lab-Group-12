docker run -d -p 8080:8080 --name root --network kadNetwork kademlia /main/start "root" "1337

for i in {1..49}
    do
        let PORT=8080+$i
        docker run -d -p ${PORT}:8080 --name "cont$i" --network mynetwork kademlia /main/start "cont$i" "${PORT}"
    done