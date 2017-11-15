<template>
  <div id="app">
    <h1>grpc web testing</h1>
    <ul>
      <li><input type="text" v-model="echoRequest"><button @click="echo">Echo</button><span>{{echoResponse}}</span></li>
    </ul>
  </div>
</template>

<script>
import { grpc, Code, Metadata } from "grpc-web-client";
import { EchoService } from "../echo/echo_pb_service";
import { EchoRequest } from "../echo/echo_pb";

const host = location.href.slice(0, -1);

export default {
  name: "app",
  data() {
    return {
      echoRequest: "",
      echoResponse: ""
    };
  },
  methods: {
    echo() {
      if (!this.echoRequest) {
        return;
      }
      const request = new EchoRequest();
      request.setMessage(this.echoRequest);
      grpc.unary(EchoService.Echo, {
        request,
        host,
        onEnd: res => {
          const { status, statusMessage, headers, message, trailers } = res;
          console.log("EchoService.Echo.onEnd.status", status, statusMessage);
          console.log("EchoService.Echo.onEnd.headers", headers);
          if (status === Code.OK && message) {
            const resp = message.toObject();
            console.log("EchoService.Echo.onEnd.message", resp);
            this.echoResponse = resp.message;
          }
          console.log("EchoService.Echo.onEnd.trailers", trailers);
        }
      });
      this.echoRequest = "";
    }
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}

h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
