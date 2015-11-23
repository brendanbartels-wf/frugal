// Autogenerated by Frugal Compiler (0.0.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library valid.src.frug_foo;

import 'dart:async';

import 'package:thrift/thrift.dart' as thrift;
import 'package:frugal/frugal.dart' as frugal;

import 'thing.dart' as t_thing;
import 'stuff.dart' as t_stuff;


const String delimiter = '.';

/// And this is a scope docstring.
class FooPublisher {
  frugal.Transport transport;
  thrift.TProtocol protocol;
  int seqId;

  FooPublisher(frugal.Provider provider) {
    var tp = provider.newTransportProtocol();
    transport = tp.transport;
    protocol = tp.protocol;
    seqId = 0;
  }

  Future publishFoo(String baz, t_thing.Thing req) {
    var op = "Foo";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    transport.preparePublish(topic);
    var oprot = protocol;
    seqId++;
    var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, seqId);
    oprot.writeMessageBegin(msg);
    req.write(oprot);
    oprot.writeMessageEnd();
    return oprot.transport.flush();
  }


<<<<<<< HEAD
  Future publishBar(String baz, t_stuff.Stuff req) {
    var op = "Bar";
=======
  /// This is an operation docstring.
  Future publishFoo(String baz, t_thing.Thing req) {
    var op = "Foo";
>>>>>>> f9437bcd253b3b469f34e48af9202b19d5a276f6
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
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


/// And this is a scope docstring.
class FooSubscriber {
  final frugal.Provider provider;

  FooSubscriber(this.provider) {}

  Future<frugal.Subscription> subscribeFoo(String baz, dynamic onThing(t_thing.Thing req)) async {
    var op = "Foo";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    var tp = provider.newTransportProtocol();
    await tp.transport.subscribe(topic);
    tp.transport.signalRead.listen((_) {
      onThing(_recvFoo(op, tp.protocol));
    });
    var sub = new frugal.Subscription(topic, tp.transport);
    tp.transport.error.listen((Error e) {;
      sub.signal(e);
    });
    return sub;
  }

  t_thing.Thing _recvFoo(String op, thrift.TProtocol iprot) {
    var tMsg = iprot.readMessageBegin();
    if (tMsg.name != op) {
      thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
      iprot.readMessageEnd();
      throw new thrift.TApplicationError(
      thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
    }
    var req = new t_thing.Thing();
    req.read(iprot);
    iprot.readMessageEnd();
    return req;
  }


<<<<<<< HEAD
  Future<frugal.Subscription> subscribeBar(String baz, dynamic onStuff(t_stuff.Stuff req)) async {
    var op = "Bar";
=======
  /// This is an operation docstring.
  Future<frugal.Subscription> subscribeFoo(String baz, dynamic onThing(t_thing.Thing req)) async {
    var op = "Foo";
>>>>>>> f9437bcd253b3b469f34e48af9202b19d5a276f6
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    var tp = provider.newTransportProtocol();
    await tp.transport.subscribe(topic);
    tp.transport.signalRead.listen((_) {
      onStuff(_recvBar(op, tp.protocol));
    });
    var sub = new frugal.Subscription(topic, tp.transport);
    tp.transport.error.listen((Error e) {;
      sub.signal(e);
    });
    return sub;
  }

  t_stuff.Stuff _recvBar(String op, thrift.TProtocol iprot) {
    var tMsg = iprot.readMessageBegin();
    if (tMsg.name != op) {
      thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
      iprot.readMessageEnd();
      throw new thrift.TApplicationError(
      thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
    }
    var req = new t_stuff.Stuff();
    req.read(iprot);
    iprot.readMessageEnd();
    return req;
  }
}

