from "xena/go-mini:1.9.2"

run "go1.9.2 download"

$repo = "github.com/PonyvilleFM/aura"

def foldercopy(dir)
  copy "#{dir}", "/root/go/src/#{$repo}/#{dir}"
end

def gobuild(pkg)
  run "go1.9.2 build #{$repo}/#{pkg} && go1.9.2 install #{$repo}/#{pkg}"
end

def cleanup()
  run "rm -rf /root/sdk /root/go/pkg /usr/local/go"
  run "apk del go1.9.2"
  flatten
end

folders = [
  "cmd",
  "internal",
  "run",
  "vendor",
]

folders.each { |x| foldercopy x }
