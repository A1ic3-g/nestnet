#!/bin/bash
docker build -t nestnet .
docker-compose up "$0"