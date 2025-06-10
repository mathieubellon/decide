port=8000 

kill:
	@lsof -ti :$(port) | xargs kill -9
	@echo "Server killed on port $(port)"

# Better reloader 
# go install github.com/mitranim/gow@latest
run:
	@echo Running on Gow -  Stop with double Ctrl-C
	@gow -e=go,html,django,css,js run .


