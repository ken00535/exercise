FROM golang:1.18-alpine AS builder
WORKDIR /server
ENV GO11MODULE=on
