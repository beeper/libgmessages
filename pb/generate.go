//go:generate protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mgoogle_messages.proto=github.com/beeper/libgmessages/pb google_messages.proto

package pb
