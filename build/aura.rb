import "build/prelude.rb"

gobuild "cmd/aura"

run %q[ apk add --virtual streamripper-deps --no-cache wget alpine-sdk glib-dev \
     && mkdir /tmp/streamripper \
     && cd /tmp/streamripper \
     && wget https://xena.greedo.xeserv.us/files/streamripper.tgz \
     && tar zxf ./streamripper.tgz \
     && cd streamripper-1.64.6 \
     && ./configure && make -j && chmod +x install-sh && make install \
     && rm -rf /tmp/streamripper \
     && apk del streamripper-deps && apk add --no-cache glib ]

cleanup
cmd "/go/src/github.com/PonyvilleFM/aura/run/aura.sh"
tag "pvfm/aura"
