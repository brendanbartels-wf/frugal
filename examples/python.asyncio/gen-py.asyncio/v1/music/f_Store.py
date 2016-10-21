#
# Autogenerated by Frugal Compiler (1.19.0)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#



import asyncio
from datetime import timedelta
import inspect

from frugal.aio.processor import FBaseProcessor
from frugal.aio.processor import FProcessorFunction
from frugal.aio.registry import FClientRegistry
from frugal.middleware import Method
from thrift.Thrift import TApplicationException
from thrift.Thrift import TMessageType
from v1.music.Store import *
from v1.music.ttypes import *


class Iface(object):
    """
    Services are the API for client and server interaction.
    Users can buy an album or enter a giveaway for a free album.
    """

    async def buyAlbum(self, ctx, ASIN, acct):
        """
        Args:
            ctx: FContext
            ASIN: string
            acct: string
        """
        pass

    async def enterAlbumGiveaway(self, ctx, email, name):
        """
        Args:
            ctx: FContext
            email: string
            name: string
        """
        pass


class Client(Iface):

    def __init__(self, transport, protocol_factory, middleware=None):
        """
        Create a new Client with a transport and protocol factory.

        Args:
            transport: FTransport
            protocol_factory: FProtocolFactory
            middleware: ServiceMiddleware or list of ServiceMiddleware
        """
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]
        transport.set_registry(FClientRegistry())
        self._transport = transport
        self._protocol_factory = protocol_factory
        self._oprot = protocol_factory.get_protocol(transport)
        self._write_lock = asyncio.Lock()
        self._methods = {
            'buyAlbum': Method(self._buyAlbum, middleware),
            'enterAlbumGiveaway': Method(self._enterAlbumGiveaway, middleware),
        }

    async def buyAlbum(self, ctx, ASIN, acct):
        """
        Args:
            ctx: FContext
            ASIN: string
            acct: string
        """
        return await self._methods['buyAlbum']([ctx, ASIN, acct])

    async def _buyAlbum(self, ctx, ASIN, acct):
        timeout = ctx.get_timeout() / 1000.0
        future = asyncio.Future()
        timed_future = asyncio.wait_for(future, timeout)
        await self._transport.register(ctx, self._recv_buyAlbum(ctx, future))
        try:
            await self._send_buyAlbum(ctx, ASIN, acct)
            result = await timed_future
        finally:
            await self._transport.unregister(ctx)
        return result

    async def _send_buyAlbum(self, ctx, ASIN, acct):
        oprot = self._oprot
        async with self._write_lock:
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('buyAlbum', TMessageType.CALL, 0)
            args = buyAlbum_args()
            args.ASIN = ASIN
            args.acct = acct
            args.write(oprot)
            oprot.writeMessageEnd()
            await oprot.get_transport().flush()

    def _recv_buyAlbum(self, ctx, future):
        def buyAlbum_callback(transport):
            iprot = self._protocol_factory.get_protocol(transport)
            iprot.read_response_headers(ctx)
            _, mtype, _ = iprot.readMessageBegin()
            if mtype == TMessageType.EXCEPTION:
                x = TApplicationException()
                x.read(iprot)
                iprot.readMessageEnd()
                future.set_exception(x)
                raise x
            result = buyAlbum_result()
            result.read(iprot)
            iprot.readMessageEnd()
            if result.error is not None:
                future.set_exception(result.error)
                return
            if result.success is not None:
                future.set_result(result.success)
                return
            x = TApplicationException(TApplicationException.MISSING_RESULT, "buyAlbum failed: unknown result")
            future.set_exception(x)
            raise x
        return buyAlbum_callback

    async def enterAlbumGiveaway(self, ctx, email, name):
        """
        Args:
            ctx: FContext
            email: string
            name: string
        """
        return await self._methods['enterAlbumGiveaway']([ctx, email, name])

    async def _enterAlbumGiveaway(self, ctx, email, name):
        timeout = ctx.get_timeout() / 1000.0
        future = asyncio.Future()
        timed_future = asyncio.wait_for(future, timeout)
        await self._transport.register(ctx, self._recv_enterAlbumGiveaway(ctx, future))
        try:
            await self._send_enterAlbumGiveaway(ctx, email, name)
            result = await timed_future
        finally:
            await self._transport.unregister(ctx)
        return result

    async def _send_enterAlbumGiveaway(self, ctx, email, name):
        oprot = self._oprot
        async with self._write_lock:
            oprot.write_request_headers(ctx)
            oprot.writeMessageBegin('enterAlbumGiveaway', TMessageType.CALL, 0)
            args = enterAlbumGiveaway_args()
            args.email = email
            args.name = name
            args.write(oprot)
            oprot.writeMessageEnd()
            await oprot.get_transport().flush()

    def _recv_enterAlbumGiveaway(self, ctx, future):
        def enterAlbumGiveaway_callback(transport):
            iprot = self._protocol_factory.get_protocol(transport)
            iprot.read_response_headers(ctx)
            _, mtype, _ = iprot.readMessageBegin()
            if mtype == TMessageType.EXCEPTION:
                x = TApplicationException()
                x.read(iprot)
                iprot.readMessageEnd()
                future.set_exception(x)
                raise x
            result = enterAlbumGiveaway_result()
            result.read(iprot)
            iprot.readMessageEnd()
            if result.success is not None:
                future.set_result(result.success)
                return
            x = TApplicationException(TApplicationException.MISSING_RESULT, "enterAlbumGiveaway failed: unknown result")
            future.set_exception(x)
            raise x
        return enterAlbumGiveaway_callback


class Processor(FBaseProcessor):

    def __init__(self, handler, middleware=None):
        """
        Create a new Processor.

        Args:
            handler: Iface
        """
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]

        super(Processor, self).__init__()
        self.add_to_processor_map('buyAlbum', _buyAlbum(Method(handler.buyAlbum, middleware), self.get_write_lock()))
        self.add_to_processor_map('enterAlbumGiveaway', _enterAlbumGiveaway(Method(handler.enterAlbumGiveaway, middleware), self.get_write_lock()))


class _buyAlbum(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._write_lock = lock

    async def process(self, ctx, iprot, oprot):
        args = buyAlbum_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = buyAlbum_result()
        try:
            ret = self._handler([ctx, args.ASIN, args.acct])
            if inspect.iscoroutine(ret):
                ret = await ret
            result.success = ret
        except PurchasingError as error:
            result.error = error
        async with self._write_lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('buyAlbum', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()


class _enterAlbumGiveaway(FProcessorFunction):

    def __init__(self, handler, lock):
        self._handler = handler
        self._write_lock = lock

    async def process(self, ctx, iprot, oprot):
        args = enterAlbumGiveaway_args()
        args.read(iprot)
        iprot.readMessageEnd()
        result = enterAlbumGiveaway_result()
        ret = self._handler([ctx, args.email, args.name])
        if inspect.iscoroutine(ret):
            ret = await ret
        result.success = ret
        async with self._write_lock:
            oprot.write_response_headers(ctx)
            oprot.writeMessageBegin('enterAlbumGiveaway', TMessageType.REPLY, 0)
            result.write(oprot)
            oprot.writeMessageEnd()
            oprot.get_transport().flush()