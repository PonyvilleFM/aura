from "xena/go-mini:1.9.1"

run "go1.9.1 download"

$repo = "github.com/PonyvilleFM/aura"

def foldercopy(dir)
  copy "#{dir}", "/root/go/src/#{$repo}/#{dir}"
end

def gobuild(pkg)
  run "go1.9.1 build #{$repo}/#{pkg} && go1.9.1 install #{$repo}/#{pkg}"
end

def cleanup()
  run "rm -rf /root/sdk /root/go/pkg /usr/local/go"
  run "apk del go1.9.1"
  flatten
end

folders = [
  "bot",
  "cmd",
  "commands",
  "doc",
  "pvfm",
  "recording",
  "run",
  "vendor",
]

folders.each { |x| foldercopy x }
