# Base Image
FROM ubuntu:22.04

LABEL org.opencontainers.image.source="https://github.com/shaharia-lab/terediX"


RUN \
    echo "Europe/Berlin" > /etc/timezone   && \
    echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections && \
    ln -sf /usr/share/zoneinfo/EET /etc/localtime

# Update system software
RUN \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    software-properties-common \
    apt-utils \
    tzdata \
    locales \
    vim \
    htop \
    zip \
    bzip2 \
    wget \
    curl \
    git \
    supervisor \
    build-essential \
    make \
    dnsutils \
    subversion \
    poppler-utils \
    fontconfig \
    xfonts-base \
    xfonts-75dpi \
    python3 \
    python3-venv \
    sudo && \
    cp -f /root/.bashrc /.bashrc && \
    mkdir -p /.config/htop/ && \
    echo "fields=0 48 17 18 38 39 40 2 46 47 49 1\nsort_key=46\nsort_direction=1\nhide_threads=0\nhide_kernel_threads=1\nhide_userland_threads=1\nshadow_other_users=0\nshow_thread_names=0\nhighlight_base_name=1\nhighlight_megabytes=1\nhighlight_threads=1\ntree_view=1\nheader_margin=1\ndetailed_cpu_time=0\ncpu_count_from_zero=0\ncolor_scheme=0\ndelay=15\nleft_meters=Hostname Clock Memory CPU Swap\nleft_meter_modes=2 2 2 1 1\nright_meters=Uptime Tasks LoadAverage AllCPUs\nright_meter_modes=2 2 2 1\n" > /.config/htop/htoprc && \
    mkdir -p /root/.config/htop/ && \
    cp -f /.config/htop/htoprc /root/.config/htop/htoprc && \
    ln -sf /usr/bin/python3 /usr/bin/python && \
    apt-get autoremove -y && \
    apt-get clean all

RUN echo 'root:pass' | chpasswd

RUN \
    locale-gen "en_US"       && \
    locale-gen "en_US.UTF-8" && \
    update-locale LANG=en_US.UTF-8 LANGUAGE=en_US

RUN adduser app --gecos "First Last,RoomNumber,WorkPhone,HomePhone" --disabled-password
RUN echo "app:pass" | chpasswd
RUN usermod -a -G adm app
RUN sudo -u app mkdir /home/app/.ssh

COPY supervisord/base_supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY supervisord/app_supervisord.conf /home/app/supervisord.conf
RUN chown app:app /home/app/supervisord.conf

CMD ["/usr/bin/supervisord"]

# Install Node.js and related packages
RUN apt-get update && \
    apt-get install -y ca-certificates curl gnupg && \
    mkdir -p /etc/apt/keyrings &&  \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg

RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list

RUN apt-get update && \
    apt-get install -y nodejs

# Install Yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor -o /usr/share/keyrings/yarnkey.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/yarnkey.gpg] https://dl.yarnpkg.com/debian stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y yarn

# Set environment variables
ENV NODE_ENV=development

## > End of Node.js

# to use GitHub repo for private Go package
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

# Download Go 1.2.2 and install it to /usr/local/go
RUN curl -s https://storage.googleapis.com/golang/go1.20.linux-amd64.tar.gz| tar -v -C /usr/local -xz

# Let's people find our Go binaries
ENV PATH $PATH:/usr/local/go/bin
RUN echo "export PATH=$PATH" > /etc/environment


## postgres
RUN \
    curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -                                        && \
    echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN \
    apt-get update                                                     && \
    apt-get install -y postgresql-14 postgresql-contrib php-pgsql && \
    rm -rf /var/lib/apt/lists/*

RUN \
    service postgresql start                                                                     && \
    sudo -u postgres psql -c "CREATE USER app WITH ENCRYPTED PASSWORD 'pass';"                   && \
    sudo -u postgres psql -c "CREATE DATABASE app ENCODING 'UTF8' OWNER app TEMPLATE template1;" && \
    sudo -u postgres psql -c "CREATE DATABASE test ENCODING 'UTF8' OWNER app TEMPLATE template1;" && \
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE app to app;"                      && \
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE test to app;"                      && \
    sudo -u postgres psql -c "ALTER USER app CREATEDB;"                                          && \
    service postgresql stop

RUN \
    echo "host    all             all             172.26.0.1/8           md5" >> /etc/postgresql/14/main/pg_hba.conf && \
    sed -i "s|#listen_addresses =.*|listen_addresses = '*'|g" /etc/postgresql/14/main/postgresql.conf

RUN \
    echo "[program:postgresql]" >> /etc/supervisor/conf.d/supervisord.conf                                                                                                                   && \
    echo "command = /usr/lib/postgresql/14/bin/postgres -D /var/lib/postgresql/14/main -c config_file=/etc/postgresql/14/main/postgresql.conf" >> /etc/supervisor/conf.d/supervisord.conf && \
    echo "user = postgres" >> /etc/supervisor/conf.d/supervisord.conf                                                                                                                        && \
    echo "autostart = true" >> /etc/supervisor/conf.d/supervisord.conf                                                                                                                       && \
    echo "stdout_events_enabled = true" >> /etc/supervisor/conf.d/supervisord.conf                                                                                                           && \
    echo "stderr_events_enabled = true" >> /etc/supervisor/conf.d/supervisord.con

# end of postgres
