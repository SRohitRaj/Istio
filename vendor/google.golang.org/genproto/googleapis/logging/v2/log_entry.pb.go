// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/logging/v2/log_entry.proto

package logging

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	_struct "github.com/golang/protobuf/ptypes/struct"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
	_type "google.golang.org/genproto/googleapis/logging/type"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// An individual entry in a log.
type LogEntry struct {
	// Required. The resource name of the log to which this log entry belongs:
	//
	//     "projects/[PROJECT_ID]/logs/[LOG_ID]"
	//     "organizations/[ORGANIZATION_ID]/logs/[LOG_ID]"
	//     "billingAccounts/[BILLING_ACCOUNT_ID]/logs/[LOG_ID]"
	//     "folders/[FOLDER_ID]/logs/[LOG_ID]"
	//
	//  A project number may optionally be used in place of PROJECT_ID. The
	//  project number is translated to its corresponding PROJECT_ID internally
	//  and the `log_name` field will contain PROJECT_ID in queries and exports.
	//
	// `[LOG_ID]` must be URL-encoded within `log_name`. Example:
	// `"organizations/1234567890/logs/cloudresourcemanager.googleapis.com%2Factivity"`.
	// `[LOG_ID]` must be less than 512 characters long and can only include the
	// following characters: upper and lower case alphanumeric characters,
	// forward-slash, underscore, hyphen, and period.
	//
	// For backward compatibility, if `log_name` begins with a forward-slash, such
	// as `/projects/...`, then the log entry is ingested as usual but the
	// forward-slash is removed. Listing the log entry will not show the leading
	// slash and filtering for a log name with a leading slash will never return
	// any results.
	LogName string `protobuf:"bytes,12,opt,name=log_name,json=logName,proto3" json:"log_name,omitempty"`
	// Required. The primary monitored resource associated with this log entry.
	// Example: a log entry that reports a database error would be
	// associated with the monitored resource designating the particular
	// database that reported the error.
	Resource *monitoredres.MonitoredResource `protobuf:"bytes,8,opt,name=resource,proto3" json:"resource,omitempty"`
	// Optional. The log entry payload, which can be one of multiple types.
	//
	// Types that are valid to be assigned to Payload:
	//	*LogEntry_ProtoPayload
	//	*LogEntry_TextPayload
	//	*LogEntry_JsonPayload
	Payload isLogEntry_Payload `protobuf_oneof:"payload"`
	// Optional. The time the event described by the log entry occurred.
	// This time is used to compute the log entry's age and to enforce
	// the logs retention period. If this field is omitted in a new log
	// entry, then Stackdriver Logging assigns it the current time.
	// Timestamps have nanosecond accuracy, but trailing zeros in the fractional
	// seconds might be omitted when the timestamp is displayed.
	//
	// Incoming log entries should have timestamps that are no more than
	// the [logs retention period](/logging/quotas) in the past,
	// and no more than 24 hours in the future. Log entries outside those time
	// boundaries will not be available when calling `entries.list`, but
	// those log entries can still be exported with
	// [LogSinks](/logging/docs/api/tasks/exporting-logs).
	Timestamp *timestamp.Timestamp `protobuf:"bytes,9,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Output only. The time the log entry was received by Stackdriver Logging.
	ReceiveTimestamp *timestamp.Timestamp `protobuf:"bytes,24,opt,name=receive_timestamp,json=receiveTimestamp,proto3" json:"receive_timestamp,omitempty"`
	// Optional. The severity of the log entry. The default value is
	// `LogSeverity.DEFAULT`.
	Severity _type.LogSeverity `protobuf:"varint,10,opt,name=severity,proto3,enum=google.logging.type.LogSeverity" json:"severity,omitempty"`
	// Optional. A unique identifier for the log entry. If you provide a value,
	// then Stackdriver Logging considers other log entries in the same project,
	// with the same `timestamp`, and with the same `insert_id` to be duplicates
	// which can be removed.  If omitted in new log entries, then Stackdriver
	// Logging assigns its own unique identifier. The `insert_id` is also used
	// to order log entries that have the same `timestamp` value.
	InsertId string `protobuf:"bytes,4,opt,name=insert_id,json=insertId,proto3" json:"insert_id,omitempty"`
	// Optional. Information about the HTTP request associated with this
	// log entry, if applicable.
	HttpRequest *_type.HttpRequest `protobuf:"bytes,7,opt,name=http_request,json=httpRequest,proto3" json:"http_request,omitempty"`
	// Optional. A set of user-defined (key, value) data that provides additional
	// information about the log entry.
	Labels map[string]string `protobuf:"bytes,11,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Output only. Additional metadata about the monitored resource.
	// Only `k8s_container`, `k8s_pod`, and `k8s_node` MonitoredResources have
	// this field populated.
	Metadata *monitoredres.MonitoredResourceMetadata `protobuf:"bytes,25,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Optional. Information about an operation associated with the log entry, if
	// applicable.
	Operation *LogEntryOperation `protobuf:"bytes,15,opt,name=operation,proto3" json:"operation,omitempty"`
	// Optional. Resource name of the trace associated with the log entry, if any.
	// If it contains a relative resource name, the name is assumed to be relative
	// to `//tracing.googleapis.com`. Example:
	// `projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824`
	Trace string `protobuf:"bytes,22,opt,name=trace,proto3" json:"trace,omitempty"`
	// Optional. The span ID within the trace associated with the log entry. For
	// Stackdriver Trace spans, this is the same format that the Stackdriver Trace
	// API v2 uses: a 16-character hexadecimal encoding of an 8-byte array, such
	// as <code>"000000000000004a"</code>.
	SpanId string `protobuf:"bytes,27,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	// Optional. Source code location information associated with the log entry,
	// if any.
	SourceLocation       *LogEntrySourceLocation `protobuf:"bytes,23,opt,name=source_location,json=sourceLocation,proto3" json:"source_location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (*LogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ba2017251165146, []int{0}
}

func (m *LogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntry.Unmarshal(m, b)
}
func (m *LogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntry.Marshal(b, m, deterministic)
}
func (m *LogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntry.Merge(m, src)
}
func (m *LogEntry) XXX_Size() int {
	return xxx_messageInfo_LogEntry.Size(m)
}
func (m *LogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntry proto.InternalMessageInfo

func (m *LogEntry) GetLogName() string {
	if m != nil {
		return m.LogName
	}
	return ""
}

func (m *LogEntry) GetResource() *monitoredres.MonitoredResource {
	if m != nil {
		return m.Resource
	}
	return nil
}

type isLogEntry_Payload interface {
	isLogEntry_Payload()
}

type LogEntry_ProtoPayload struct {
	ProtoPayload *any.Any `protobuf:"bytes,2,opt,name=proto_payload,json=protoPayload,proto3,oneof"`
}

type LogEntry_TextPayload struct {
	TextPayload string `protobuf:"bytes,3,opt,name=text_payload,json=textPayload,proto3,oneof"`
}

type LogEntry_JsonPayload struct {
	JsonPayload *_struct.Struct `protobuf:"bytes,6,opt,name=json_payload,json=jsonPayload,proto3,oneof"`
}

func (*LogEntry_ProtoPayload) isLogEntry_Payload() {}

func (*LogEntry_TextPayload) isLogEntry_Payload() {}

func (*LogEntry_JsonPayload) isLogEntry_Payload() {}

func (m *LogEntry) GetPayload() isLogEntry_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *LogEntry) GetProtoPayload() *any.Any {
	if x, ok := m.GetPayload().(*LogEntry_ProtoPayload); ok {
		return x.ProtoPayload
	}
	return nil
}

func (m *LogEntry) GetTextPayload() string {
	if x, ok := m.GetPayload().(*LogEntry_TextPayload); ok {
		return x.TextPayload
	}
	return ""
}

func (m *LogEntry) GetJsonPayload() *_struct.Struct {
	if x, ok := m.GetPayload().(*LogEntry_JsonPayload); ok {
		return x.JsonPayload
	}
	return nil
}

func (m *LogEntry) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *LogEntry) GetReceiveTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.ReceiveTimestamp
	}
	return nil
}

func (m *LogEntry) GetSeverity() _type.LogSeverity {
	if m != nil {
		return m.Severity
	}
	return _type.LogSeverity_DEFAULT
}

func (m *LogEntry) GetInsertId() string {
	if m != nil {
		return m.InsertId
	}
	return ""
}

func (m *LogEntry) GetHttpRequest() *_type.HttpRequest {
	if m != nil {
		return m.HttpRequest
	}
	return nil
}

func (m *LogEntry) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *LogEntry) GetMetadata() *monitoredres.MonitoredResourceMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *LogEntry) GetOperation() *LogEntryOperation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *LogEntry) GetTrace() string {
	if m != nil {
		return m.Trace
	}
	return ""
}

func (m *LogEntry) GetSpanId() string {
	if m != nil {
		return m.SpanId
	}
	return ""
}

func (m *LogEntry) GetSourceLocation() *LogEntrySourceLocation {
	if m != nil {
		return m.SourceLocation
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LogEntry) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LogEntry_OneofMarshaler, _LogEntry_OneofUnmarshaler, _LogEntry_OneofSizer, []interface{}{
		(*LogEntry_ProtoPayload)(nil),
		(*LogEntry_TextPayload)(nil),
		(*LogEntry_JsonPayload)(nil),
	}
}

func _LogEntry_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LogEntry)
	// payload
	switch x := m.Payload.(type) {
	case *LogEntry_ProtoPayload:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ProtoPayload); err != nil {
			return err
		}
	case *LogEntry_TextPayload:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.TextPayload)
	case *LogEntry_JsonPayload:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.JsonPayload); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LogEntry.Payload has unexpected type %T", x)
	}
	return nil
}

func _LogEntry_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LogEntry)
	switch tag {
	case 2: // payload.proto_payload
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(any.Any)
		err := b.DecodeMessage(msg)
		m.Payload = &LogEntry_ProtoPayload{msg}
		return true, err
	case 3: // payload.text_payload
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Payload = &LogEntry_TextPayload{x}
		return true, err
	case 6: // payload.json_payload
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(_struct.Struct)
		err := b.DecodeMessage(msg)
		m.Payload = &LogEntry_JsonPayload{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LogEntry_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LogEntry)
	// payload
	switch x := m.Payload.(type) {
	case *LogEntry_ProtoPayload:
		s := proto.Size(x.ProtoPayload)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LogEntry_TextPayload:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.TextPayload)))
		n += len(x.TextPayload)
	case *LogEntry_JsonPayload:
		s := proto.Size(x.JsonPayload)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Additional information about a potentially long-running operation with which
// a log entry is associated.
type LogEntryOperation struct {
	// Optional. An arbitrary operation identifier. Log entries with the
	// same identifier are assumed to be part of the same operation.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Optional. An arbitrary producer identifier. The combination of
	// `id` and `producer` must be globally unique.  Examples for `producer`:
	// `"MyDivision.MyBigCompany.com"`, `"github.com/MyProject/MyApplication"`.
	Producer string `protobuf:"bytes,2,opt,name=producer,proto3" json:"producer,omitempty"`
	// Optional. Set this to True if this is the first log entry in the operation.
	First bool `protobuf:"varint,3,opt,name=first,proto3" json:"first,omitempty"`
	// Optional. Set this to True if this is the last log entry in the operation.
	Last                 bool     `protobuf:"varint,4,opt,name=last,proto3" json:"last,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntryOperation) Reset()         { *m = LogEntryOperation{} }
func (m *LogEntryOperation) String() string { return proto.CompactTextString(m) }
func (*LogEntryOperation) ProtoMessage()    {}
func (*LogEntryOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ba2017251165146, []int{1}
}

func (m *LogEntryOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntryOperation.Unmarshal(m, b)
}
func (m *LogEntryOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntryOperation.Marshal(b, m, deterministic)
}
func (m *LogEntryOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntryOperation.Merge(m, src)
}
func (m *LogEntryOperation) XXX_Size() int {
	return xxx_messageInfo_LogEntryOperation.Size(m)
}
func (m *LogEntryOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntryOperation.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntryOperation proto.InternalMessageInfo

func (m *LogEntryOperation) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *LogEntryOperation) GetProducer() string {
	if m != nil {
		return m.Producer
	}
	return ""
}

func (m *LogEntryOperation) GetFirst() bool {
	if m != nil {
		return m.First
	}
	return false
}

func (m *LogEntryOperation) GetLast() bool {
	if m != nil {
		return m.Last
	}
	return false
}

// Additional information about the source code location that produced the log
// entry.
type LogEntrySourceLocation struct {
	// Optional. Source file name. Depending on the runtime environment, this
	// might be a simple name or a fully-qualified name.
	File string `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	// Optional. Line within the source file. 1-based; 0 indicates no line number
	// available.
	Line int64 `protobuf:"varint,2,opt,name=line,proto3" json:"line,omitempty"`
	// Optional. Human-readable name of the function or method being invoked, with
	// optional context such as the class or package name. This information may be
	// used in contexts such as the logs viewer, where a file and line number are
	// less meaningful. The format can vary by language. For example:
	// `qual.if.ied.Class.method` (Java), `dir/package.func` (Go), `function`
	// (Python).
	Function             string   `protobuf:"bytes,3,opt,name=function,proto3" json:"function,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntrySourceLocation) Reset()         { *m = LogEntrySourceLocation{} }
func (m *LogEntrySourceLocation) String() string { return proto.CompactTextString(m) }
func (*LogEntrySourceLocation) ProtoMessage()    {}
func (*LogEntrySourceLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ba2017251165146, []int{2}
}

func (m *LogEntrySourceLocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntrySourceLocation.Unmarshal(m, b)
}
func (m *LogEntrySourceLocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntrySourceLocation.Marshal(b, m, deterministic)
}
func (m *LogEntrySourceLocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntrySourceLocation.Merge(m, src)
}
func (m *LogEntrySourceLocation) XXX_Size() int {
	return xxx_messageInfo_LogEntrySourceLocation.Size(m)
}
func (m *LogEntrySourceLocation) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntrySourceLocation.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntrySourceLocation proto.InternalMessageInfo

func (m *LogEntrySourceLocation) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func (m *LogEntrySourceLocation) GetLine() int64 {
	if m != nil {
		return m.Line
	}
	return 0
}

func (m *LogEntrySourceLocation) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "google.logging.v2.LogEntry")
	proto.RegisterMapType((map[string]string)(nil), "google.logging.v2.LogEntry.LabelsEntry")
	proto.RegisterType((*LogEntryOperation)(nil), "google.logging.v2.LogEntryOperation")
	proto.RegisterType((*LogEntrySourceLocation)(nil), "google.logging.v2.LogEntrySourceLocation")
}

func init() { proto.RegisterFile("google/logging/v2/log_entry.proto", fileDescriptor_8ba2017251165146) }

var fileDescriptor_8ba2017251165146 = []byte{
	// 748 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x5d, 0x6f, 0xf3, 0x34,
	0x14, 0x7e, 0xd3, 0xbe, 0xb4, 0xa9, 0xdb, 0x77, 0x1f, 0xd6, 0x58, 0xb3, 0x6e, 0x88, 0xb2, 0xf1,
	0x51, 0x6e, 0x52, 0xa9, 0xdc, 0x6c, 0x6c, 0x12, 0x5a, 0x27, 0xb4, 0x4d, 0xea, 0x60, 0xf2, 0xd0,
	0x2e, 0x50, 0xa5, 0xca, 0x4b, 0xdc, 0xcc, 0x90, 0xda, 0xc1, 0x71, 0x2a, 0xfa, 0x53, 0xf8, 0x0b,
	0x5c, 0xf3, 0x2b, 0xf8, 0x45, 0x5c, 0x22, 0x9f, 0x38, 0x69, 0x69, 0xa7, 0x72, 0x77, 0x8e, 0xcf,
	0xf3, 0x9c, 0xe7, 0xc9, 0xf1, 0x71, 0xd0, 0x67, 0x91, 0x94, 0x51, 0xcc, 0xfa, 0xb1, 0x8c, 0x22,
	0x2e, 0xa2, 0xfe, 0x7c, 0x60, 0xc2, 0x09, 0x13, 0x5a, 0x2d, 0xfc, 0x44, 0x49, 0x2d, 0xf1, 0x7e,
	0x0e, 0xf1, 0x2d, 0xc4, 0x9f, 0x0f, 0x3a, 0x27, 0x96, 0x45, 0x13, 0xde, 0xa7, 0x42, 0x48, 0x4d,
	0x35, 0x97, 0x22, 0xcd, 0x09, 0x9d, 0xb3, 0x95, 0xea, 0x4c, 0x0a, 0xae, 0xa5, 0x62, 0xe1, 0x44,
	0xb1, 0x54, 0x66, 0x2a, 0x60, 0x16, 0xf4, 0xe5, 0x9a, 0xb0, 0x5e, 0x24, 0xac, 0xff, 0xaa, 0x75,
	0x32, 0x51, 0xec, 0xb7, 0x8c, 0xa5, 0x7a, 0x1b, 0xce, 0x58, 0x4c, 0xd9, 0x9c, 0x29, 0xae, 0xad,
	0xcb, 0xce, 0x91, 0xc5, 0x41, 0xf6, 0x92, 0x4d, 0xfb, 0x54, 0x14, 0xa5, 0x93, 0xf5, 0x52, 0xaa,
	0x55, 0x16, 0x14, 0x02, 0x9f, 0xae, 0x57, 0x35, 0x9f, 0xb1, 0x54, 0xd3, 0x59, 0x92, 0x03, 0x4e,
	0xff, 0xaa, 0x23, 0x77, 0x24, 0xa3, 0xef, 0xcd, 0x48, 0xf0, 0x11, 0x72, 0x8d, 0xb8, 0xa0, 0x33,
	0xe6, 0xb5, 0xba, 0x4e, 0xaf, 0x41, 0xea, 0xb1, 0x8c, 0x7e, 0xa0, 0x33, 0x86, 0x2f, 0x90, 0x5b,
	0x7c, 0xa3, 0xe7, 0x76, 0x9d, 0x5e, 0x73, 0xf0, 0x89, 0x6f, 0x47, 0x47, 0x13, 0xee, 0x3f, 0x14,
	0x93, 0x20, 0x16, 0x44, 0x4a, 0x38, 0xbe, 0x44, 0x1f, 0x40, 0x6b, 0x92, 0xd0, 0x45, 0x2c, 0x69,
	0xe8, 0x55, 0x80, 0x7f, 0x50, 0xf0, 0x0b, 0x6f, 0xfe, 0xb5, 0x58, 0xdc, 0xbd, 0x23, 0x2d, 0xc8,
	0x1f, 0x73, 0x2c, 0x3e, 0x43, 0x2d, 0xcd, 0x7e, 0xd7, 0x25, 0xb7, 0x6a, 0x6c, 0xdd, 0xbd, 0x23,
	0x4d, 0x73, 0x5a, 0x80, 0xae, 0x50, 0xeb, 0x97, 0x54, 0x8a, 0x12, 0x54, 0x03, 0x81, 0xf6, 0x86,
	0xc0, 0x13, 0x8c, 0xc6, 0xb0, 0x0d, 0xbc, 0x60, 0x9f, 0xa3, 0x46, 0x39, 0x15, 0xaf, 0x01, 0xd4,
	0xce, 0x06, 0xf5, 0xa7, 0x02, 0x41, 0x96, 0x60, 0x7c, 0x8b, 0xf6, 0x15, 0x0b, 0x18, 0x9f, 0xb3,
	0xc9, 0xb2, 0x83, 0xf7, 0xbf, 0x1d, 0xf6, 0x2c, 0xa9, 0x3c, 0xc1, 0x57, 0xc8, 0x2d, 0x6e, 0xdc,
	0x43, 0x5d, 0xa7, 0xb7, 0x33, 0xe8, 0xfa, 0x6b, 0x8b, 0x69, 0x56, 0xc3, 0x1f, 0xc9, 0xe8, 0xc9,
	0xe2, 0x48, 0xc9, 0xc0, 0xc7, 0xa8, 0xc1, 0x45, 0xca, 0x94, 0x9e, 0xf0, 0xd0, 0x7b, 0x0f, 0xf7,
	0xe6, 0xe6, 0x07, 0xf7, 0x21, 0xbe, 0x41, 0xad, 0xd5, 0xc5, 0xf3, 0xea, 0x60, 0xef, 0xed, 0xf6,
	0x77, 0x5a, 0x27, 0x24, 0xc7, 0x91, 0xe6, 0xeb, 0x32, 0xc1, 0xdf, 0xa1, 0x5a, 0x4c, 0x5f, 0x58,
	0x9c, 0x7a, 0xcd, 0x6e, 0xb5, 0xd7, 0x1c, 0x7c, 0xe5, 0x6f, 0x3c, 0x1b, 0xbf, 0xd8, 0x22, 0x7f,
	0x04, 0x48, 0x88, 0x89, 0xa5, 0xe1, 0x6b, 0xe4, 0xce, 0x98, 0xa6, 0x21, 0xd5, 0xd4, 0x3b, 0x02,
	0x07, 0x5f, 0x6c, 0x5d, 0x9f, 0x07, 0x0b, 0x26, 0x25, 0x0d, 0x0f, 0x51, 0x43, 0x26, 0x4c, 0xc1,
	0x63, 0xf4, 0x76, 0xa1, 0xc7, 0xe7, 0x5b, 0x6c, 0xfc, 0x58, 0x60, 0xc9, 0x92, 0x86, 0x0f, 0xd0,
	0x47, 0x5a, 0xd1, 0x80, 0x79, 0x87, 0x30, 0xa5, 0x3c, 0xc1, 0x6d, 0x54, 0x4f, 0x13, 0x2a, 0xcc,
	0xf4, 0x8e, 0xe1, 0xbc, 0x66, 0xd2, 0xfb, 0x10, 0x13, 0xb4, 0x9b, 0xdb, 0x99, 0xc4, 0x32, 0xc8,
	0x85, 0xdb, 0x20, 0xfc, 0xf5, 0x16, 0xe1, 0x27, 0x60, 0x8c, 0x2c, 0x81, 0xec, 0xa4, 0xff, 0xc9,
	0x3b, 0x17, 0xa8, 0xb9, 0x32, 0x20, 0xbc, 0x87, 0xaa, 0xbf, 0xb2, 0x85, 0xe7, 0x80, 0xae, 0x09,
	0x8d, 0xc7, 0x39, 0x8d, 0x33, 0x06, 0xcf, 0xa4, 0x41, 0xf2, 0xe4, 0xdb, 0xca, 0xb9, 0x33, 0x6c,
	0xa0, 0xba, 0xdd, 0xf0, 0x53, 0x8e, 0xf6, 0x37, 0x3e, 0x14, 0xef, 0xa0, 0x0a, 0x0f, 0x6d, 0xab,
	0x0a, 0x0f, 0x71, 0x07, 0xb9, 0x89, 0x92, 0x61, 0x16, 0x30, 0x65, 0x9b, 0x95, 0xb9, 0x51, 0x99,
	0x72, 0x95, 0x6a, 0x78, 0x50, 0x2e, 0xc9, 0x13, 0x8c, 0xd1, 0xfb, 0x98, 0xa6, 0x1a, 0x96, 0xc8,
	0x25, 0x10, 0x9f, 0x8e, 0xd1, 0xe1, 0xdb, 0x9f, 0x66, 0xd0, 0x53, 0x1e, 0x33, 0xab, 0x08, 0x31,
	0x74, 0xe0, 0x22, 0x37, 0x5f, 0x25, 0x10, 0x1b, 0x1f, 0xd3, 0x4c, 0x04, 0x30, 0xbf, 0x6a, 0xee,
	0xa3, 0xc8, 0x87, 0x7f, 0x38, 0xe8, 0xe3, 0x40, 0xce, 0x36, 0xe7, 0x39, 0xfc, 0x50, 0xa8, 0x3e,
	0xc2, 0xff, 0xc0, 0xf9, 0xf9, 0xdc, 0x62, 0x22, 0x19, 0x53, 0x11, 0xf9, 0x52, 0x45, 0xfd, 0x88,
	0x09, 0x78, 0x5f, 0xfd, 0xbc, 0x44, 0x13, 0x9e, 0xae, 0xfc, 0xec, 0x2f, 0x6d, 0xf8, 0x8f, 0xe3,
	0xfc, 0x59, 0x69, 0xdf, 0xe6, 0xec, 0x9b, 0x58, 0x66, 0xa1, 0xb9, 0x2c, 0xd0, 0x79, 0x1e, 0xfc,
	0x5d, 0x54, 0xc6, 0x50, 0x19, 0xdb, 0xca, 0xf8, 0x79, 0xf0, 0x52, 0x83, 0xde, 0xdf, 0xfc, 0x1b,
	0x00, 0x00, 0xff, 0xff, 0xc7, 0x91, 0xad, 0xed, 0x47, 0x06, 0x00, 0x00,
}
