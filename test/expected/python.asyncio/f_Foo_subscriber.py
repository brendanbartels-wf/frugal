#
# Autogenerated by Frugal Compiler (1.19.2)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#



import inspect
import sys
import traceback

from thrift.Thrift import TApplicationException
from thrift.Thrift import TMessageType
from thrift.Thrift import TType
from frugal.middleware import Method
from frugal.subscription import FSubscription

from valid.ttypes import *




class FooSubscriber(object):
    """
    And this is a scope docstring.
    """

    _DELIMITER = '.'

    def __init__(self, provider, middleware=None):
        """
        Create a new FooSubscriber.

        Args:
            provider: FScopeProvider
            middleware: ServiceMiddleware or list of ServiceMiddleware
        """

        if middleware and not isinstance(middleware, list):
            middleware = [middleware]
        self._middleware = middleware
        self._transport, self._protocol_factory = provider.new()

    async def subscribe_Foo(self, baz, Foo_handler):
        """
        This is an operation docstring.
        
        Args:
            baz: string
            Foo_handler: function which takes FContext and Thing
        """

        op = 'Foo'
        prefix = 'foo.bar.{}.qux.'.format(baz)
        topic = '{}Foo{}{}'.format(prefix, self._DELIMITER, op)

        await self._transport.subscribe(topic, self._recv_Foo(self._protocol_factory, op, Foo_handler))

    def _recv_Foo(self, protocol_factory, op, handler):
        method = Method(handler, self._middleware)

        async def callback(transport):
            iprot = protocol_factory.get_protocol(transport)
            ctx = iprot.read_request_headers()
            mname, _, _ = iprot.readMessageBegin()
            if mname != op:
                iprot.skip(TType.STRUCT)
                iprot.readMessageEnd()
                raise TApplicationException(TApplicationException.UNKNOWN_METHOD)
            req = Thing()
            req.read(iprot)
            iprot.readMessageEnd()
            try:
                ret = method([ctx, req])
                if inspect.iscoroutine(ret):
                    await ret
            except:
                traceback.print_exc()
                sys.exit(1)

        return callback



    async def subscribe_Bar(self, baz, Bar_handler):
        """
        Args:
            baz: string
            Bar_handler: function which takes FContext and Stuff
        """

        op = 'Bar'
        prefix = 'foo.bar.{}.qux.'.format(baz)
        topic = '{}Foo{}{}'.format(prefix, self._DELIMITER, op)

        await self._transport.subscribe(topic, self._recv_Bar(self._protocol_factory, op, Bar_handler))

    def _recv_Bar(self, protocol_factory, op, handler):
        method = Method(handler, self._middleware)

        async def callback(transport):
            iprot = protocol_factory.get_protocol(transport)
            ctx = iprot.read_request_headers()
            mname, _, _ = iprot.readMessageBegin()
            if mname != op:
                iprot.skip(TType.STRUCT)
                iprot.readMessageEnd()
                raise TApplicationException(TApplicationException.UNKNOWN_METHOD)
            req = Stuff()
            req.read(iprot)
            iprot.readMessageEnd()
            try:
                ret = method([ctx, req])
                if inspect.iscoroutine(ret):
                    await ret
            except:
                traceback.print_exc()
                sys.exit(1)

        return callback




