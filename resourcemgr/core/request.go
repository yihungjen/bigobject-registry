package core

type ResourceReq struct {
	Action string
	Name   string
	Resp   chan *ResourceInfo
}

func NewResourceReq() *ResourceReq {
	return &ResourceReq{Resp: make(chan *ResourceInfo, 1)}
}

type ResourceInfo struct {
	Err error
}
