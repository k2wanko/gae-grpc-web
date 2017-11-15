// package: net.k2lab.test.grpc.testing.echo
// file: echo.proto

var jspb = require("google-protobuf");
var echo_pb = require("./echo_pb");
var EchoService = {
  serviceName: "net.k2lab.test.grpc.testing.echo.EchoService"
};
EchoService.Echo = {
  methodName: "Echo",
  service: EchoService,
  requestStream: false,
  responseStream: false,
  requestType: echo_pb.EchoRequest,
  responseType: echo_pb.EchoResponse
};
module.exports = {
  EchoService: EchoService,
};

