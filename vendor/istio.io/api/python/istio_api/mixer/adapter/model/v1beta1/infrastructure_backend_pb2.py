# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: mixer/adapter/model/v1beta1/infrastructure_backend.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from google.rpc import status_pb2 as google_dot_rpc_dot_status__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='mixer/adapter/model/v1beta1/infrastructure_backend.proto',
  package='istio.mixer.adapter.model.v1beta1',
  syntax='proto3',
  serialized_pb=_b('\n8mixer/adapter/model/v1beta1/infrastructure_backend.proto\x12!istio.mixer.adapter.model.v1beta1\x1a\x19google/protobuf/any.proto\x1a\x17google/rpc/status.proto\"\xf4\x01\n\x14\x43reateSessionRequest\x12,\n\x0e\x61\x64\x61pter_config\x18\x01 \x01(\x0b\x32\x14.google.protobuf.Any\x12\x62\n\x0einferred_types\x18\x02 \x03(\x0b\x32J.istio.mixer.adapter.model.v1beta1.CreateSessionRequest.InferredTypesEntry\x1aJ\n\x12InferredTypesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12#\n\x05value\x18\x02 \x01(\x0b\x32\x14.google.protobuf.Any:\x02\x38\x01\"O\n\x15\x43reateSessionResponse\x12\x12\n\nsession_id\x18\x01 \x01(\t\x12\"\n\x06status\x18\x02 \x01(\x0b\x32\x12.google.rpc.Status\"\xea\x01\n\x0fValidateRequest\x12,\n\x0e\x61\x64\x61pter_config\x18\x01 \x01(\x0b\x32\x14.google.protobuf.Any\x12]\n\x0einferred_types\x18\x02 \x03(\x0b\x32\x45.istio.mixer.adapter.model.v1beta1.ValidateRequest.InferredTypesEntry\x1aJ\n\x12InferredTypesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12#\n\x05value\x18\x02 \x01(\x0b\x32\x14.google.protobuf.Any:\x02\x38\x01\"6\n\x10ValidateResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status\")\n\x13\x43loseSessionRequest\x12\x12\n\nsession_id\x18\x01 \x01(\t\":\n\x14\x43loseSessionResponse\x12\"\n\x06status\x18\x01 \x01(\x0b\x32\x12.google.rpc.Status2\x92\x03\n\x15InfrastructureBackend\x12s\n\x08Validate\x12\x32.istio.mixer.adapter.model.v1beta1.ValidateRequest\x1a\x33.istio.mixer.adapter.model.v1beta1.ValidateResponse\x12\x82\x01\n\rCreateSession\x12\x37.istio.mixer.adapter.model.v1beta1.CreateSessionRequest\x1a\x38.istio.mixer.adapter.model.v1beta1.CreateSessionResponse\x12\x7f\n\x0c\x43loseSession\x12\x36.istio.mixer.adapter.model.v1beta1.CloseSessionRequest\x1a\x37.istio.mixer.adapter.model.v1beta1.CloseSessionResponseB-Z(istio.io/api/mixer/adapter/model/v1beta1\x80\x01\x01\x62\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_any__pb2.DESCRIPTOR,google_dot_rpc_dot_status__pb2.DESCRIPTOR,])




_CREATESESSIONREQUEST_INFERREDTYPESENTRY = _descriptor.Descriptor(
  name='InferredTypesEntry',
  full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest.InferredTypesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest.InferredTypesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest.InferredTypesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=_descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001')),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=318,
  serialized_end=392,
)

_CREATESESSIONREQUEST = _descriptor.Descriptor(
  name='CreateSessionRequest',
  full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='adapter_config', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest.adapter_config', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='inferred_types', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionRequest.inferred_types', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_CREATESESSIONREQUEST_INFERREDTYPESENTRY, ],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=148,
  serialized_end=392,
)


_CREATESESSIONRESPONSE = _descriptor.Descriptor(
  name='CreateSessionResponse',
  full_name='istio.mixer.adapter.model.v1beta1.CreateSessionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='session_id', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionResponse.session_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='status', full_name='istio.mixer.adapter.model.v1beta1.CreateSessionResponse.status', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=394,
  serialized_end=473,
)


_VALIDATEREQUEST_INFERREDTYPESENTRY = _descriptor.Descriptor(
  name='InferredTypesEntry',
  full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest.InferredTypesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest.InferredTypesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest.InferredTypesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=_descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001')),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=318,
  serialized_end=392,
)

_VALIDATEREQUEST = _descriptor.Descriptor(
  name='ValidateRequest',
  full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='adapter_config', full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest.adapter_config', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='inferred_types', full_name='istio.mixer.adapter.model.v1beta1.ValidateRequest.inferred_types', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_VALIDATEREQUEST_INFERREDTYPESENTRY, ],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=476,
  serialized_end=710,
)


_VALIDATERESPONSE = _descriptor.Descriptor(
  name='ValidateResponse',
  full_name='istio.mixer.adapter.model.v1beta1.ValidateResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='istio.mixer.adapter.model.v1beta1.ValidateResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=712,
  serialized_end=766,
)


_CLOSESESSIONREQUEST = _descriptor.Descriptor(
  name='CloseSessionRequest',
  full_name='istio.mixer.adapter.model.v1beta1.CloseSessionRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='session_id', full_name='istio.mixer.adapter.model.v1beta1.CloseSessionRequest.session_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=768,
  serialized_end=809,
)


_CLOSESESSIONRESPONSE = _descriptor.Descriptor(
  name='CloseSessionResponse',
  full_name='istio.mixer.adapter.model.v1beta1.CloseSessionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='status', full_name='istio.mixer.adapter.model.v1beta1.CloseSessionResponse.status', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=811,
  serialized_end=869,
)

_CREATESESSIONREQUEST_INFERREDTYPESENTRY.fields_by_name['value'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_CREATESESSIONREQUEST_INFERREDTYPESENTRY.containing_type = _CREATESESSIONREQUEST
_CREATESESSIONREQUEST.fields_by_name['adapter_config'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_CREATESESSIONREQUEST.fields_by_name['inferred_types'].message_type = _CREATESESSIONREQUEST_INFERREDTYPESENTRY
_CREATESESSIONRESPONSE.fields_by_name['status'].message_type = google_dot_rpc_dot_status__pb2._STATUS
_VALIDATEREQUEST_INFERREDTYPESENTRY.fields_by_name['value'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_VALIDATEREQUEST_INFERREDTYPESENTRY.containing_type = _VALIDATEREQUEST
_VALIDATEREQUEST.fields_by_name['adapter_config'].message_type = google_dot_protobuf_dot_any__pb2._ANY
_VALIDATEREQUEST.fields_by_name['inferred_types'].message_type = _VALIDATEREQUEST_INFERREDTYPESENTRY
_VALIDATERESPONSE.fields_by_name['status'].message_type = google_dot_rpc_dot_status__pb2._STATUS
_CLOSESESSIONRESPONSE.fields_by_name['status'].message_type = google_dot_rpc_dot_status__pb2._STATUS
DESCRIPTOR.message_types_by_name['CreateSessionRequest'] = _CREATESESSIONREQUEST
DESCRIPTOR.message_types_by_name['CreateSessionResponse'] = _CREATESESSIONRESPONSE
DESCRIPTOR.message_types_by_name['ValidateRequest'] = _VALIDATEREQUEST
DESCRIPTOR.message_types_by_name['ValidateResponse'] = _VALIDATERESPONSE
DESCRIPTOR.message_types_by_name['CloseSessionRequest'] = _CLOSESESSIONREQUEST
DESCRIPTOR.message_types_by_name['CloseSessionResponse'] = _CLOSESESSIONRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

CreateSessionRequest = _reflection.GeneratedProtocolMessageType('CreateSessionRequest', (_message.Message,), dict(

  InferredTypesEntry = _reflection.GeneratedProtocolMessageType('InferredTypesEntry', (_message.Message,), dict(
    DESCRIPTOR = _CREATESESSIONREQUEST_INFERREDTYPESENTRY,
    __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
    # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.CreateSessionRequest.InferredTypesEntry)
    ))
  ,
  DESCRIPTOR = _CREATESESSIONREQUEST,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.CreateSessionRequest)
  ))
_sym_db.RegisterMessage(CreateSessionRequest)
_sym_db.RegisterMessage(CreateSessionRequest.InferredTypesEntry)

CreateSessionResponse = _reflection.GeneratedProtocolMessageType('CreateSessionResponse', (_message.Message,), dict(
  DESCRIPTOR = _CREATESESSIONRESPONSE,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.CreateSessionResponse)
  ))
_sym_db.RegisterMessage(CreateSessionResponse)

ValidateRequest = _reflection.GeneratedProtocolMessageType('ValidateRequest', (_message.Message,), dict(

  InferredTypesEntry = _reflection.GeneratedProtocolMessageType('InferredTypesEntry', (_message.Message,), dict(
    DESCRIPTOR = _VALIDATEREQUEST_INFERREDTYPESENTRY,
    __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
    # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.ValidateRequest.InferredTypesEntry)
    ))
  ,
  DESCRIPTOR = _VALIDATEREQUEST,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.ValidateRequest)
  ))
_sym_db.RegisterMessage(ValidateRequest)
_sym_db.RegisterMessage(ValidateRequest.InferredTypesEntry)

ValidateResponse = _reflection.GeneratedProtocolMessageType('ValidateResponse', (_message.Message,), dict(
  DESCRIPTOR = _VALIDATERESPONSE,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.ValidateResponse)
  ))
_sym_db.RegisterMessage(ValidateResponse)

CloseSessionRequest = _reflection.GeneratedProtocolMessageType('CloseSessionRequest', (_message.Message,), dict(
  DESCRIPTOR = _CLOSESESSIONREQUEST,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.CloseSessionRequest)
  ))
_sym_db.RegisterMessage(CloseSessionRequest)

CloseSessionResponse = _reflection.GeneratedProtocolMessageType('CloseSessionResponse', (_message.Message,), dict(
  DESCRIPTOR = _CLOSESESSIONRESPONSE,
  __module__ = 'mixer.adapter.model.v1beta1.infrastructure_backend_pb2'
  # @@protoc_insertion_point(class_scope:istio.mixer.adapter.model.v1beta1.CloseSessionResponse)
  ))
_sym_db.RegisterMessage(CloseSessionResponse)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('Z(istio.io/api/mixer/adapter/model/v1beta1\200\001\001'))
_CREATESESSIONREQUEST_INFERREDTYPESENTRY.has_options = True
_CREATESESSIONREQUEST_INFERREDTYPESENTRY._options = _descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001'))
_VALIDATEREQUEST_INFERREDTYPESENTRY.has_options = True
_VALIDATEREQUEST_INFERREDTYPESENTRY._options = _descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001'))

_INFRASTRUCTUREBACKEND = _descriptor.ServiceDescriptor(
  name='InfrastructureBackend',
  full_name='istio.mixer.adapter.model.v1beta1.InfrastructureBackend',
  file=DESCRIPTOR,
  index=0,
  options=None,
  serialized_start=872,
  serialized_end=1274,
  methods=[
  _descriptor.MethodDescriptor(
    name='Validate',
    full_name='istio.mixer.adapter.model.v1beta1.InfrastructureBackend.Validate',
    index=0,
    containing_service=None,
    input_type=_VALIDATEREQUEST,
    output_type=_VALIDATERESPONSE,
    options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CreateSession',
    full_name='istio.mixer.adapter.model.v1beta1.InfrastructureBackend.CreateSession',
    index=1,
    containing_service=None,
    input_type=_CREATESESSIONREQUEST,
    output_type=_CREATESESSIONRESPONSE,
    options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CloseSession',
    full_name='istio.mixer.adapter.model.v1beta1.InfrastructureBackend.CloseSession',
    index=2,
    containing_service=None,
    input_type=_CLOSESESSIONREQUEST,
    output_type=_CLOSESESSIONRESPONSE,
    options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_INFRASTRUCTUREBACKEND)

DESCRIPTOR.services_by_name['InfrastructureBackend'] = _INFRASTRUCTUREBACKEND

# @@protoc_insertion_point(module_scope)
