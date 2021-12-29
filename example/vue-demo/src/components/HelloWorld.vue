<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <p>
      For a guide and recipes on how to configure / customize this project,<br>
      check out the
      <a href="https://cli.vuejs.org" target="_blank" rel="noopener">vue-cli documentation</a>.
    </p>
    <h3>Installed CLI Plugins</h3>
    <a @click="goAbout()">GO TO About</a>
  </div>
</template>

<script>
import prefixUrl from 'prefix-uri/prefix-url'

export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },
  methods: {
    goAbout() {
      var baseUrl = window.location.origin + prefixUrl('/about')
      console.log("baseUrl: ", baseUrl)
      this.$router.push(baseUrl)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
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
