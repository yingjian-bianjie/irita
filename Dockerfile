FROM lenking/wenchangchain_builder:rocskdb_v6.15.5

WORKDIR /irita

COPY . .

RUN WITH_ROCKSDB=yes make build \
    && mv build/irita /usr/local/bin

