FROM golang:1.12-stretch

ARG code_dir

ENV ENV development
ENV CODE_DIR $code_dir
COPY . $code_dir
WORKDIR $code_dir
RUN go build -o story-book-api
CMD "./story-book-api"
#CMD "tail -f /dev/null"