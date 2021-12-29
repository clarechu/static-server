<template>
  <div id="app">
    <h1>{{ msg }}</h1>
    <h3>HTTP Client</h3>
    <button @click="health()"> Check Health </button>
    <p>
      URL TO <a href="{{baseUrl}}" target="_blank" rel="noopener">{{baseUrl}}</a>.
    </p>
    <p>
      <h3> response : {{response}}</h3>
    </p>
  </div>
</template>

<script>
import prefixUrl from '../utils/prefix-url'

export default {
  name: "app",
  props: {
    msg: String
  },
  data() {
    return {
      baseUrl: '',
      response: '',
    }
  },
  created: function () {
    this.baseUrl = window.location.origin
  },
  methods: {
    health() {
      var baseUrl = window.location.origin + prefixUrl('/health')
      this.baseUrl = baseUrl
      console.log("baseUrl: ", baseUrl)
      this.axios.get(baseUrl).then((response) => {
        console.log(response.code)
        this.response = response
      }).catch((error) => {
        console.log(error)
        this.response = error
      })
    }
  }
};
</script>