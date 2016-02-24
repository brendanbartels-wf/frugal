part of frugal;

var _encoder = new Utf8Encoder();
var _decoder = new Utf8Decoder();

/**
 * This is an internal-only class. Don't use it!
 */
class Headers {
  static const _V0 = 0x00;

  /// Encode the headers
  static Uint8List encode(Map<String, String> headers) {
    var size = 0;
    // Get total frame size headers
    if (headers != null && headers.length > 0) {
      for (var name in headers.keys) {
        // 4 bytes each for name, value length
        size += 8 + name.length + headers[name].length;
      }
    }

    // Header buff = [version (1 byte), size (4 bytes), headers (size bytes)]
    var buff = new Uint8List(5 + size);

    // Write version
    buff[0] = _V0;

    // Write size
    _writeInt(size, buff, 1);

    // Write headers
    if (headers != null && headers.length > 0) {
      var i = 5;
      for (var name in headers.keys) {
        // Write name length
        _writeInt(name.length, buff, i);
        i += 4;
        // Write name
        _writeString(name, buff, i);
        i += name.length;

        // Write value length
        var value = headers[name];
        _writeInt(value.length, buff, i);
        i += 4;
        _writeString(value, buff, i);
        i += value.length;
      }
    }
    return buff;
  }

  /// Reads the headers from a TTransport
  static Map<String, String> read(TTransport transport) {
    // Buffer version
    var buff = new Uint8List(5);
    transport.readAll(buff, 0, 1);

    // Support more versions when available
    if (buff[0] != _V0) {
      throw new FProtocolError(TProtocolErrorType.BAD_VERSION, "unsupported header version ${buff[0]}");
    }

    // Read size
    transport.readAll(buff, 1, 4);
    var size = _readInt(buff, 1);

    // Read the rest of the header bytes into a buffer
    buff = new Uint8List(size);
    transport.readAll(buff, 0, size);

    return _readPairs(buff, 0, size);
  }

  /// Returns the headers from Frugal frame
  static Map<String, String> decodeFromFrame(Uint8List frame) {
    if (frame.length < 5) {
      throw new FProtocolError(TProtocolErrorType.INVALID_DATA, "invalid frame size ${frame.length}");
    }

    // Support more versions when available
    if (frame[0] != _V0) {
      throw new FProtocolError(TProtocolErrorType.BAD_VERSION, "unsupported header version ${frame[0]}");
    }

    return _readPairs(frame, 5, _readInt(frame, 1) + 5);
  }

  static Map<String, String> _readPairs(Uint8List buff, int start, int end) {
    var headers = {};
    for (var i = start; i < end; i) {
      // Read header name
      var nameSize = _readInt(buff, i);
      i += 4;
      if (i > end || i + nameSize > end) {
        throw new FProtocolError(TProtocolErrorType.INVALID_DATA, "invalid protocol header name");
      }
      var name = _decoder.convert(buff, i, i + nameSize);
      i += nameSize;

      // Read header value
      var valueSize = _readInt(buff, i);
      i += 4;
      if (i > end || i + valueSize > end) {
        throw new FProtocolError(TProtocolErrorType.INVALID_DATA, "invalid protocol header value");
      }
      var value = _decoder.convert(buff, i, i + valueSize);
      i += valueSize;

      // Set the pair
      headers[name] = value;
    }
    return headers;
  }

  static int _readInt(Uint8List buff, int i) {
    return ((buff[i] & 0xff) << 24) |
        ((buff[i + 1] & 0xff) << 16) |
        ((buff[i + 2] & 0xff) << 8) |
        (buff[i + 3] & 0xff);
  }

  static void _writeInt(int i, Uint8List buff, int i1) {
    buff[i1] = (0xff & (i >> 24));
    buff[i1 + 1] = (0xff & (i >> 16));
    buff[i1 + 2] = (0xff & (i >> 8));
    buff[i1 + 3] = (0xff & (i));
  }

  static void _writeString(String s, Uint8List buff, int i) {
    var strBytes = _encoder.convert(s);
    buff.setRange(i, i + strBytes.length, strBytes);
  }
}
