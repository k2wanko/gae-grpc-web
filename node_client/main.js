const messages = require('../echo/echo_pb')
const services = require('../echo/echo_pb_service')

const {grpc, Code} = require('grpc-web-client')

const host = 'https://gae-grpc-web.appspot.com:443'

function main() {
  const request = new messages.EchoRequest()
  request.setMessage('Hello')
  grpc.unary(services.EchoService.Echo, {
    request,
    host,
    onEnd: res => {
      const { status, statusMessage, headers, message, trailers } = res
      console.log("EchoService.Echo.onEnd.status", status, statusMessage)
      console.log("EchoService.Echo.onEnd.headers", headers)
      if (status === Code.OK && message) {
        const resp = message.toObject()
        console.log("EchoService.Echo.onEnd.message", resp)
      }
      console.log("EchoService.Echo.onEnd.trailers", trailers)
    }
  })
}

main()
