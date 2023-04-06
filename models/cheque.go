package models

import "database/sql"

type NullChequeLogo struct {
	Image  sql.NullString  `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	Left   sql.NullFloat64 `protobuf:"fixed32,2,opt,name=left,proto3" json:"left,omitempty"`
	Right  sql.NullFloat64 `protobuf:"fixed32,3,opt,name=right,proto3" json:"right,omitempty"`
	Top    sql.NullFloat64 `protobuf:"fixed32,4,opt,name=top,proto3" json:"top,omitempty"`
	Bottom sql.NullFloat64 `protobuf:"fixed32,5,opt,name=bottom,proto3" json:"bottom,omitempty"`
}
