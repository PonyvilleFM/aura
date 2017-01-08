from "xena/go:latest"

$repo = "github.com/PonyvilleFM/aura"

def foldercopy(dir)
  copy "#{dir}", "/go/src/#{$repo}/#{dir}"
end

def gobuild(pkg)
  run "go build #{$repo}/#{pkg} && go install #{$repo}/#{pkg}"
end

def cleanup()
  run "rm -rf /go/pkg /usr/local/go"
  run "apk del go"
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
  "vendor-log"
]

folders.each { |x| foldercopy x }
