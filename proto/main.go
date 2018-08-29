package main

import (
	"fmt"

	proto "github.com/gogo/protobuf/proto"
)

type TaskMeta struct {
	MaxConcurrency             int32             `protobuf:"varint,1,opt,name=max_concurrency,json=maxConcurrency,proto3" json:"max_concurrency,omitempty"`
	LastCompletedTimestampUnix int64             `protobuf:"varint,2,opt,name=last_completed_timestamp_unix,json=lastCompletedTimestampUnix,proto3" json:"last_completed_timestamp_unix,omitempty"`
	CurrentlyRunning           []*TaskMetaRunner `protobuf:"bytes,3,rep,name=currently_running,json=currentlyRunning" json:"currently_running,omitempty"`
}

func (m *TaskMeta) Reset()         { *m = TaskMeta{} }
func (m *TaskMeta) String() string { return proto.CompactTextString(m) }
func (*TaskMeta) ProtoMessage()    {}

type TaskMetaRunner struct {
	NowTimestampUnix int64  `protobuf:"varint,1,opt,name=now_timestamp_unix,json=nowTimestampUnix,proto3" json:"now_timestamp_unix,omitempty"`
	Try              uint32 `protobuf:"varint,2,opt,name=try,proto3" json:"try,omitempty"`
	RunID            uint64 `protobuf:"varint,3,opt,name=run_id,json=runId,proto3" json:"run_id,omitempty"`
}

func (m *TaskMetaRunner) Reset()         { *m = TaskMetaRunner{} }
func (m *TaskMetaRunner) String() string { return proto.CompactTextString(m) }
func (*TaskMetaRunner) ProtoMessage()    {}

func main() {
	t := TaskMeta{
		MaxConcurrency:             1,
		LastCompletedTimestampUnix: 388,
		CurrentlyRunning: []*TaskMetaRunner{
			{
				NowTimestampUnix: 123,
				Try:              1,
				RunID:            3242,
			},
		},
	}

	bytes, err := proto.Marshal(&t)
	fmt.Printf("bytes: %s\nlen: %d\n err: %v\n", bytes, len(bytes), err)
	nt := TaskMeta{}

	err = proto.Unmarshal(bytes, &nt)
	fmt.Printf("nt: %#v\n err: %v\n", nt, err)

}
