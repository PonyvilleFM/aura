import "build/prelude.rb"

env "GO111MODULE" => "on"

gobuild "cmd/aerial"

cleanup

cmd "/aerial"
tag "pvfm/aerial"
