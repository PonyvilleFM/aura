from "xena/go:1.9.4"

$repo = "github.com/PonyvilleFM/aura"

def foldercopy(dir)
  copy "#{dir}", "/root/go/src/#{$repo}/#{dir}"
end

def gobuild(pkg)
  run "go1.9.4 build #{$repo}/#{pkg} && go1.9.4 install #{$repo}/#{pkg}"
end

def cleanup()
  run "rm -rf /root/sdk /root/go/pkg /usr/local/go /usr/local/bin/go"
  flatten
end

folders = [
  "cmd",
  "internal",
  "run",
  "vendor",
]

folders.each { |x| foldercopy x }
