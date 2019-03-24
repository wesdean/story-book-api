# story-book-api
A platform for authoring interactive stories

**Example Docker build**

`docker build --build-arg code_dir=/go/src/github.com/wesdean/story-book-api -t story-book-api . && docker run -p 3000:3000 --name story-book-api story-book-api`