.PHONY: client server

server:
	@docker build -t chat-server ./server
	@docker run --rm --network="host" -p 8080:8080 --name chat-server chat-server

client:
	@docker build -t chat-client ./client
	@docker run --rm --network="host" -it chat-client