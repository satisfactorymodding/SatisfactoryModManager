FROM ubuntu:22.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# Pull webkit2gtk-4.1 build dependencies
RUN sed -i 's/^#\s*deb-src/deb-src/' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends build-essential curl ca-certificates && \
    apt-get build-dep -y libwebkit2gtk-4.1-dev

# Build freetype 2.13
ARG FREETYPE_VERSION=2.13.3
RUN mkdir -p /build/freetype && cd /build/freetype && \
    curl -sSL "https://download.savannah.gnu.org/releases/freetype/freetype-${FREETYPE_VERSION}.tar.xz" | tar -xJ --strip-components=1 && \
    ./configure --prefix=/usr/local --enable-shared --disable-static && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig && \
    mkdir -p /etc/dpkg && echo "libfreetype 6 libfreetype6" >> /etc/dpkg/shlibs.override

# Build webkit2gtk-4.1
RUN mkdir -p /build/webkit && cd /build/webkit && \
    apt-get source webkit2gtk && \
    cd webkit2gtk-*/ && \
    sed -i 's/ENABLE_SOUP2=YES/ENABLE_SOUP2=NO/' debian/rules && \
    sed -i 's/ENABLE_GTK4=YES/ENABLE_GTK4=NO/' debian/rules && \
    DEB_BUILD_OPTIONS="nocheck nodoc parallel=$(nproc)" dpkg-buildpackage -b -uc -us

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Used in the cache key of actions/setup-go, would be "undefined" if missing
ENV ImageOS=ubuntu22

# Install packages normally installed on GHA
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        sudo \
        build-essential \
        pkg-config \
        file

RUN curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && \
    sudo apt-get install -y nodejs

RUN (type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y)) \
	&& sudo mkdir -p -m 755 /etc/apt/keyrings \
	&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
	&& cat $out | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
	&& sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
	&& sudo mkdir -p -m 755 /etc/apt/sources.list.d \
	&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
	&& sudo apt update \
	&& sudo apt install gh -y

# Install our build of webkit2gtk-4.1
COPY --from=builder /build/webkit/*.deb /tmp/debs/

RUN apt-get install -y --no-install-recommends /tmp/debs/*.deb && \
    rm -rf /tmp/debs
