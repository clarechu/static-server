<template>
  <div class="navbar">
    <h1>This is an navbar page</h1>
    <Navbar msg="Welcome to Your Vue.js App to Navbar page"/>
  </div>
</template>

<script>
// @ is an alias to /src
import Navbar from '@/components/Navbar.vue'

export default {
  name: 'Home',
  components: {
    Navbar
  }
}
</script>
