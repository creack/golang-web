NAME            = gotest
IMAGE           = creack/golang-web
GOFILES         = $(shell find . -name '*.go') tpl static

all             : build

build           : .docker_id

.docker_id      : Makefile Dockerfile $(GOFILES)
		docker build -t $(NAME) .
		docker inspect -f '{{.Id}}' $(NAME) > $@

run             : build
		docker run -itP $(NAME)

push            : release
		docker push $(IMAGE)

clean           :
		rm -f .docker_id .release

re              : clean all

.release        : Dockerfile.release .docker_id
		docker run --rm --entrypoint /bin/sh $(NAME) -c 'tar cf - /$(NAME) /etc/ssl' > $@ || (rm -f $@; false)
		docker build --rm -t $(IMAGE) -f Dockerfile.release . || (rm -f $@; false)

release         : .release

.PHONY          : all build run push clean re release
