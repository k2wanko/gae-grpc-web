const {EchoHistoryRequest} = require('../echo/echo_pb')
const {EchoService} = require('../echo/echo_pb_service')

const {grpc, Code} = require('grpc-web-client')

const host = 'https://gae-grpc-web.appspot.com:443'

function main() {
  const request = new EchoHistoryRequest()
  request.setLimit(10)
  grpc.invoke(EchoService.EchoHistory, {
    request,
    host,
    onMessage: msg => {
      const resp = msg.toObject()
      console.log("fetchEchoHistory", "onMessage", resp)
    },
    onEnd: () => {}
  })
}

main()
