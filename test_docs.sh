docker run -i -d \
    --name teredix-dev \
    -v "$HOME"/.ssh:/home/app/.ssh \
    -v "$HOME"/Projects/terediX:/home/app/src \
    -p 3000:3000 \
    hf-dev-ubuntu-postgres-node:16