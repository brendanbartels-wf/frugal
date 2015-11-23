// Autogenerated by Frugal Compiler (0.0.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library event.src.frug_events;

import 'dart:async';

import 'package:thrift/thrift.dart' as thrift;
import 'package:frugal/frugal.dart' as frugal;

import 'event.dart' as t_event;


const String delimiter = '.';

/// This docstring gets added to the generated code because it has
/// the @ sign.
class EventsPublisher {
  frugal.Transport transport;
  thrift.TProtocol protocol;
  int seqId;

  EventsPublisher(frugal.Provider provider) {
    var tp = provider.newTransportProtocol();
    transport = tp.transport;
    protocol = tp.protocol;
    seqId = 0;
  }

  /// This is a docstring.
  Future publishEventCreated(String user, t_event.Event req) {
    var op = "EventCreated";
    var prefix = "foo.${user}.";
    var topic = "${prefix}Events${delimiter}${op}";
    transport.preparePublish(topic);
    var oprot = protocol;
    seqId++;
    var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, seqId);
    oprot.writeMessageBegin(msg);
    req.write(oprot);
    oprot.writeMessageEnd();
    return oprot.transport.flush();
  }
}


/// This docstring gets added to the generated code because it has
/// the @ sign.
class EventsSubscriber {
  final frugal.Provider provider;

  EventsSubscriber(this.provider) {}

  /// This is a docstring.
  Future<frugal.Subscription> subscribeEventCreated(String user, dynamic onEvent(t_event.Event req)) async {
    var op = "EventCreated";
    var prefix = "foo.${user}.";
    var topic = "${prefix}Events${delimiter}${op}";
    var tp = provider.newTransportProtocol();
    await tp.transport.subscribe(topic);
    tp.transport.signalRead.listen((_) {
      onEvent(_recvEventCreated(op, tp.protocol));
    });
    var sub = new frugal.Subscription(topic, tp.transport);
    tp.transport.error.listen((Error e) {;
      sub.signal(e);
    });
    return sub;
  }

  t_event.Event _recvEventCreated(String op, thrift.TProtocol iprot) {
    var tMsg = iprot.readMessageBegin();
    if (tMsg.name != op) {
      thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
      iprot.readMessageEnd();
      throw new thrift.TApplicationError(
      thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
    }
    var req = new t_event.Event();
    req.read(iprot);
    iprot.readMessageEnd();
    return req;
  }
}

