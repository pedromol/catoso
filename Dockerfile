FROM lscr.io/linuxserver/ffmpeg:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
            tzdata git build-essential cmake pkg-config wget unzip libgtk2.0-dev \
            curl ca-certificates libcurl4-openssl-dev libssl-dev \
            libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev \
            libjpeg-turbo8-dev libpng-dev libtiff-dev libdc1394-dev nasm && \
            rm -rf /var/lib/apt/lists/*

ARG OPENCV_VERSION="4.8.1"
ENV OPENCV_VERSION $OPENCV_VERSION

ARG OPENCV_FILE="https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip"
ENV OPENCV_FILE $OPENCV_FILE

ARG OPENCV_CONTRIB_FILE="https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip"
ENV OPENCV_CONTRIB_FILE $OPENCV_CONTRIB_FILE

RUN curl -Lo opencv.zip ${OPENCV_FILE} && \
            unzip -q opencv.zip && \
            curl -Lo opencv_contrib.zip ${OPENCV_CONTRIB_FILE} && \
            unzip -q opencv_contrib.zip && \
            rm opencv.zip opencv_contrib.zip && \
            cd opencv-${OPENCV_VERSION} && \
            mkdir build

WORKDIR /opencv-${OPENCV_VERSION}/build

RUN [ "$(lscpu | head -n 1 | awk '{print $2}')" = aarch64 ] && export EXTRA_FLAGS='-D ENABLE_NEON=ON' || export EXTRA_FLAGS='-D ENABLE_NEON=OFF' && \
    cmake $EXTRA_FLAGS -D CMAKE_BUILD_TYPE=RELEASE \
            -D WITH_IPP=OFF \
            -D WITH_OPENGL=OFF \
            -D WITH_QT=OFF \
            -D CMAKE_INSTALL_PREFIX=/usr/local \
            -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
            -D OPENCV_ENABLE_NONFREE=ON \
            -D WITH_JASPER=OFF \
            -D WITH_TBB=ON \
            -D BUILD_JPEG=ON \
            -D WITH_SIMD=ON \
            -D ENABLE_LIBJPEG_TURBO_SIMD=ON \
            -D BUILD_DOCS=OFF \
            -D BUILD_EXAMPLES=OFF \
            -D BUILD_TESTS=OFF \
            -D BUILD_PERF_TESTS=ON \
            -D BUILD_opencv_java=NO \
            -D BUILD_opencv_python=NO \
            -D BUILD_opencv_python2=NO \
            -D BUILD_opencv_python3=NO \
            -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
    make -j $(nproc --all) && \
    make preinstall && make install && ldconfig && \
    cd / && rm -rf opencv*

ENV ASDF_DIR=/root/.asdf
WORKDIR /go/src/app

RUN git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.13.1 && \
    echo '. "$HOME/.asdf/asdf.sh"' >> ~/.bashrc && \
    . "$HOME/.asdf/asdf.sh" && \
    asdf plugin add golang https://github.com/asdf-community/asdf-golang.git && \
    asdf install golang 1.21.3 && \
    asdf global golang 1.21.3
    
COPY cmd /go/src/app/cmd
COPY pkg /go/src/app/pkg
COPY go.mod /go/src/app/
COPY go.sum /go/src/app/
COPY data/haarcascade_frontalcatface_extended.xml /
COPY data/haarcascade_frontalcatface.xml /

RUN . "$HOME/.asdf/asdf.sh" && \
    go build -o catoso cmd/catoso/main.go && \
    cp catoso /usr/local/bin

WORKDIR /
RUN rm -Rf /go

ENTRYPOINT ["catoso"]
