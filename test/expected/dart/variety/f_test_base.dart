// Autogenerated by Frugal Compiler (1.7.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library variety.src.f_test_base;

import 'dart:typed_data' show Uint8List;
import 'package:thrift/thrift.dart';
import 'package:variety/variety.dart' as t_variety;
import 'package:actual_base/actual_base.dart' as t_actual_base;

class TestBase implements TBase {
  static final TStruct _STRUCT_DESC = new TStruct("TestBase");
  static final TField _BASE_STRUCT_FIELD_DESC = new TField("base_struct", TType.STRUCT, 1);

  t_actual_base.thing _base_struct;
  static const int BASE_STRUCT = 1;


  TestBase() {
  }

  t_actual_base.thing get base_struct => this._base_struct;

  set base_struct(t_actual_base.thing base_struct) {
    this._base_struct = base_struct;
  }

  bool isSetBase_struct() => this.base_struct != null;

  unsetBase_struct() {
    this.base_struct = null;
  }

  getFieldValue(int fieldID) {
    switch (fieldID) {
      case BASE_STRUCT:
        return this.base_struct;
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  setFieldValue(int fieldID, Object value) {
    switch(fieldID) {
      case BASE_STRUCT:
        if(value == null) {
          unsetBase_struct();
        } else {
          this.base_struct = value;
        }
        break;

      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  // Returns true if the field corresponding to fieldID is set (has been assigned a value) and false otherwise
  bool isSet(int fieldID) {
    switch(fieldID) {
      case BASE_STRUCT:
        return isSetBase_struct();
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  read(TProtocol iprot) {
    TField field;
    iprot.readStructBegin();
    while(true) {
      field = iprot.readFieldBegin();
      if(field.type == TType.STOP) {
        break;
      }
      switch(field.id) {
        case BASE_STRUCT:
          if(field.type == TType.STRUCT) {
            base_struct = new t_actual_base.thing();
            base_struct.read(iprot);
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        default:
          TProtocolUtil.skip(iprot, field.type);
          break;
      }
      iprot.readFieldEnd();
    }
    iprot.readStructEnd();

    // check for required fields of primitive type, which can't be checked in the validate method
    validate();
  }

  write(TProtocol oprot) {
    validate();

    oprot.writeStructBegin(_STRUCT_DESC);
    if(this.base_struct != null) {
      oprot.writeFieldBegin(_BASE_STRUCT_FIELD_DESC);
      base_struct.write(oprot);
      oprot.writeFieldEnd();
    }
    oprot.writeFieldStop();
    oprot.writeStructEnd();
  }

  String toString() {
    StringBuffer ret = new StringBuffer("TestBase(");

    ret.write("base_struct:");
    if(this.base_struct == null) {
      ret.write("null");
    } else {
      ret.write(this.base_struct);
    }

    ret.write(")");

    return ret.toString();
  }

  validate() {
    // check for required fields
    // check that fields of type enum have valid values
  }
}
