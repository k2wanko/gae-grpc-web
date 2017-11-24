<template>
  <div id="app">
    <h1>grpc web testing</h1>
    <ul>
      <li><input type="text" v-model="echoRequest"><button @click="echo">Echo</button><span>{{echoResponse}}</span></li>
    </ul>
    <h2>History</h2>
    <ul>
      <li v-for="echo in history" :key="echo.id">{{echo.message}}</li>
    </ul>
  </div>
</template>

<script>
import { grpc, Code, Metadata } from "grpc-web-client";
import { EchoService } from "../echo/echo_pb_service";
import { EchoRequest, EchoHistoryRequest } from "../echo/echo_pb";

const host = location.protocol + "//" + location.hostname + ":" + location.port;

export default {
  name: "app",
  data() {
    return {
      echoRequest: "",
      echoResponse: "",
      history: []
    };
  },
  mounted() {
    this.fetchEchoHistory();
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
            this.echoResponse = resp.echo.message;
            if (this.history.length >= 10) {
              this.history.pop();
            }
            this.history.unshift(resp.echo);
          }
          console.log("EchoService.Echo.onEnd.trailers", trailers);

          this.echoRequest = "";
        }
      });
    },
    fetchEchoHistory() {
      const request = new EchoHistoryRequest();
      request.setLimit(10);
      grpc.invoke(EchoService.EchoHistory, {
        request,
        host,
        onMessage: msg => {
          const resp = msg.toObject();
          console.log("fetchEchoHistory", "onMessage", resp);
          this.history.push(resp.echo);
        },
        onEnd: () => {}
      });
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
  text-align: center;
}

h1,
h2 {
  font-weight: normal;
}

ul {
  /* list-style-type: none; */
  padding: 0;
}

/* li {
  display: inline-block;
  margin: 0 10px;
} */

a {
  color: #42b983;
}
</style>
