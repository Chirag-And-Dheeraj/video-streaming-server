clean:
	@if [ -d "video" ]; then rm -r video; fi
	@if [ -d "segments" ]; then rm -r segments; fi
	@if [ -f "database.db" ]; then rm database.db; fi
	@echo Clean up complete.

init:
	@mkdir video
	@mkdir segments
	@echo Initalized.

start:
	go run main.go

cleanstart:
	@echo Clean starting...
	make clean init start
