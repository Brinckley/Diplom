kafka_up:
	@docker-compose up init-kafka

producer_up:
	@docker-compose up --build producer

consul_up:
	@docker-compose up -d consul

consul_register:
	@./consul.d/register-services.sh

up:
	consul_up consul_register kafka_up producer_up