redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

redis-ping:
	docker exec -it redis redis-cli ping

.PHONY: redis redis-ping