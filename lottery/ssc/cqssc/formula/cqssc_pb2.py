# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: cqssc.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='cqssc.proto',
  package='cqssc',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x0b\x63qssc.proto\x12\x05\x63qssc\"(\n\x0bValidateReq\x12\x0b\n\x03tag\x18\x01 \x01(\t\x12\x0c\n\x04\x63ode\x18\x02 \x03(\t\"\x1b\n\x0cValidateResp\x12\x0b\n\x03num\x18\x01 \x01(\x05\"7\n\nComputeReq\x12\x0b\n\x03tag\x18\x01 \x01(\t\x12\x0c\n\x04\x63ode\x18\x02 \x03(\t\x12\x0e\n\x06result\x18\x03 \x03(\t\"\x1a\n\x0b\x43omputeResp\x12\x0b\n\x03num\x18\x01 \x01(\x05\x32t\n\x07\x46ormula\x12\x35\n\x08Validate\x12\x12.cqssc.ValidateReq\x1a\x13.cqssc.ValidateResp\"\x00\x12\x32\n\x07\x43ompute\x12\x11.cqssc.ComputeReq\x1a\x12.cqssc.ComputeResp\"\x00\x62\x06proto3')
)




_VALIDATEREQ = _descriptor.Descriptor(
  name='ValidateReq',
  full_name='cqssc.ValidateReq',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='tag', full_name='cqssc.ValidateReq.tag', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='code', full_name='cqssc.ValidateReq.code', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=22,
  serialized_end=62,
)


_VALIDATERESP = _descriptor.Descriptor(
  name='ValidateResp',
  full_name='cqssc.ValidateResp',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='num', full_name='cqssc.ValidateResp.num', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=64,
  serialized_end=91,
)


_COMPUTEREQ = _descriptor.Descriptor(
  name='ComputeReq',
  full_name='cqssc.ComputeReq',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='tag', full_name='cqssc.ComputeReq.tag', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='code', full_name='cqssc.ComputeReq.code', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='result', full_name='cqssc.ComputeReq.result', index=2,
      number=3, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=93,
  serialized_end=148,
)


_COMPUTERESP = _descriptor.Descriptor(
  name='ComputeResp',
  full_name='cqssc.ComputeResp',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='num', full_name='cqssc.ComputeResp.num', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=150,
  serialized_end=176,
)

DESCRIPTOR.message_types_by_name['ValidateReq'] = _VALIDATEREQ
DESCRIPTOR.message_types_by_name['ValidateResp'] = _VALIDATERESP
DESCRIPTOR.message_types_by_name['ComputeReq'] = _COMPUTEREQ
DESCRIPTOR.message_types_by_name['ComputeResp'] = _COMPUTERESP
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ValidateReq = _reflection.GeneratedProtocolMessageType('ValidateReq', (_message.Message,), dict(
  DESCRIPTOR = _VALIDATEREQ,
  __module__ = 'cqssc_pb2'
  # @@protoc_insertion_point(class_scope:cqssc.ValidateReq)
  ))
_sym_db.RegisterMessage(ValidateReq)

ValidateResp = _reflection.GeneratedProtocolMessageType('ValidateResp', (_message.Message,), dict(
  DESCRIPTOR = _VALIDATERESP,
  __module__ = 'cqssc_pb2'
  # @@protoc_insertion_point(class_scope:cqssc.ValidateResp)
  ))
_sym_db.RegisterMessage(ValidateResp)

ComputeReq = _reflection.GeneratedProtocolMessageType('ComputeReq', (_message.Message,), dict(
  DESCRIPTOR = _COMPUTEREQ,
  __module__ = 'cqssc_pb2'
  # @@protoc_insertion_point(class_scope:cqssc.ComputeReq)
  ))
_sym_db.RegisterMessage(ComputeReq)

ComputeResp = _reflection.GeneratedProtocolMessageType('ComputeResp', (_message.Message,), dict(
  DESCRIPTOR = _COMPUTERESP,
  __module__ = 'cqssc_pb2'
  # @@protoc_insertion_point(class_scope:cqssc.ComputeResp)
  ))
_sym_db.RegisterMessage(ComputeResp)



_FORMULA = _descriptor.ServiceDescriptor(
  name='Formula',
  full_name='cqssc.Formula',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=178,
  serialized_end=294,
  methods=[
  _descriptor.MethodDescriptor(
    name='Validate',
    full_name='cqssc.Formula.Validate',
    index=0,
    containing_service=None,
    input_type=_VALIDATEREQ,
    output_type=_VALIDATERESP,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Compute',
    full_name='cqssc.Formula.Compute',
    index=1,
    containing_service=None,
    input_type=_COMPUTEREQ,
    output_type=_COMPUTERESP,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_FORMULA)

DESCRIPTOR.services_by_name['Formula'] = _FORMULA

# @@protoc_insertion_point(module_scope)